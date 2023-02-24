package businessLogic

import (
	"backend/utils"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

func NewAccountService(accountDb utils.AccountDatabase) *AccountService {
	return &AccountService{
		db: accountDb,
	}
}

func (as AccountService) Close() {
	as.db.Close()
}

func (as AccountService) SignIn(email string, password string) (utils.AccountApiResponse, error) {
	account, err := as.db.GetByEmail(email)

	if err != nil {
		return utils.AccountApiResponse{}, errors.New("invalidEmail")
	}

	if ok, err := as.validatePassword(password, account.PasswordHash); err != nil || !ok {
		return utils.AccountApiResponse{}, errors.New("invalidPassword")
	}

	return account.OmitPasswordHash(), nil
}

func (as AccountService) SignUp(email string, password string, name string, address string, postcode string) (utils.AccountApiResponse, error) {
	account, err := as.db.Add(utils.Account{
		PasswordHash: as.HashPassword(password),
		ShippingDetails: utils.ShippingDetails{
			Name:     name,
			Email:    email,
			Postcode: postcode,
			Address:  address,
		},
	})

	if err != nil {
		return utils.AccountApiResponse{}, errors.New("Failed to sign up user")
	}

	return account.OmitPasswordHash(), nil
}

func (as AccountService) Update(account utils.UpdateAccount) (utils.AccountApiResponse, error) {
	oldAccount, err := as.db.GetById(account.Id)

	if err != nil {
		return utils.AccountApiResponse{}, err
	}

	updateAccount := as.populateSetAccountValues(oldAccount, account)

	updatedAccount, err := as.db.Update(updateAccount)

	if err != nil {
		return utils.AccountApiResponse{}, errors.New("Failed to update account")
	}

	return updatedAccount.OmitPasswordHash(), nil
}

func (as AccountService) populateSetAccountValues(oldAccount utils.Account, newAccount utils.UpdateAccount) utils.Account {
	if newAccount.Address != "" {
		oldAccount.Address = newAccount.Address
	}

	if newAccount.Email != "" {
		oldAccount.Email = newAccount.Email
	}

	if newAccount.Name != "" {
		oldAccount.Name = newAccount.Name
	}

	if newAccount.Postcode != "" {
		oldAccount.Postcode = newAccount.Postcode
	}

	if newAccount.Password != "" {
		oldAccount.PasswordHash = as.HashPassword(newAccount.Password)
	}

	return oldAccount
}

func (as AccountService) GetById(accountId int) (utils.AccountApiResponse, error) {
	account, err := as.db.GetById(accountId)

	if err != nil {
		return utils.AccountApiResponse{}, errors.New("Failed to get account with id: " + strconv.Itoa(accountId))
	}

	return account.OmitPasswordHash(), nil
}

type AccountService struct {
	db utils.AccountDatabase
}

func (as AccountService) validatePassword(password string, hashNsalt string) (bool, error) {
	hash, salt, ok := strings.Cut(hashNsalt, ":")

	if !ok {
		return false, errors.New("Failed to split hash and salt")
	}

	newHash := as.StrongishHash(password, salt)

	return hash == newHash, nil
}

// Don't try this at home kids. Use a cloud service for authentication stuff.
func (as AccountService) StrongishHash(password string, salt string) string {
	return base64.StdEncoding.EncodeToString(pbkdf2.Key([]byte(password), []byte(salt), 1000, 64, sha512.New))
}

func (as AccountService) HashPassword(password string) string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)

	if err != nil {
		panic("Failed to generate random string: " + err.Error())
	}

	salt := base64.StdEncoding.EncodeToString(randomBytes)
	hash := as.StrongishHash(password, string(salt))

	return hash + ":" + salt
}
