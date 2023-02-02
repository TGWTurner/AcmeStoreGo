package businessLogic

import (
	"bjssStoreGo/backend/utils"
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

func (as AccountService) SignIn(email string, password string) (utils.Account, error) {
	account, err := as.db.GetByEmail(email)

	if err != nil {
		return utils.Account{}, errors.New("Failed to get account for email: " + email)
	}

	if ok, err := as.validatePassword(password, account.PasswordHash); err != nil || !ok {
		return utils.Account{}, errors.New("Invalid password for account with email: " + email)
	}

	return account, nil
}

func (as AccountService) SignUp(email string, password string, name string, address string, postcode string) (utils.Account, error) {
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
		return utils.Account{}, errors.New("Failed to sign up user")
	}

	return account, nil
}

func (as AccountService) Update(account utils.Account, newPassword ...string) (utils.Account, error) {
	if len(newPassword) > 0 {
		account.PasswordHash = as.HashPassword(account.PasswordHash)
	}

	updatedAccount, err := as.db.Update(account)

	if err != nil {
		return utils.Account{}, errors.New("Failed to update account with accountId: " + strconv.Itoa(account.Id))
	}

	return updatedAccount, nil
}

func (as AccountService) GetById(accountId int) (utils.Account, error) {
	account, err := as.db.GetById(accountId)

	if err != nil {
		return utils.Account{}, errors.New("Failed to get account with id: " + strconv.Itoa(accountId))
	}

	return account, nil
}

type AccountService struct {
	db utils.AccountDatabase
}

func (as AccountService) validatePassword(password string, hashNsalt string) (bool, error) {
	hash, salt, ok := strings.Cut(hashNsalt, ":")

	if !ok {
		return false, errors.New("Failed to split has and salt")
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
