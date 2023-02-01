package businessLogic

import (
	"bjssStoreGo/backend/utils"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

func NewAccountService(accountDb utils.AccountDatabase) *AccountService {
	return &AccountService{
		accountDb: accountDb,
	}
}

func (as AccountService) signIn(email string, password string) (utils.Account, error) {
	account, err := as.accountDb.GetByEmail(email)

	if err != nil {
		return utils.Account{}, errors.New("Failed to get account for email: " + email)
	}

	if ok, err := validatePassword(password, account.PasswordHash); err != nil || !ok {
		return utils.Account{}, errors.New("Invalid password for account with email: " + email)
	}

	return account, nil
}

func (as AccountService) signUp(email string, password string, name string, address string, postcode string) (utils.Account, error) {
	account, err := as.accountDb.Add(utils.Account{
		PasswordHash: hashPassword(password),
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

func (as AccountService) update(account utils.Account, newPassword ...string) (utils.Account, error) {
	if len(newPassword) > 0 {
		account.PasswordHash = hashPassword(account.PasswordHash)
	}

	updatedAccount, err := as.accountDb.Update(account)

	if err != nil {
		return utils.Account{}, errors.New("Failed to update account with accountId: " + strconv.Itoa(account.Id))
	}

	return updatedAccount, nil
}

func (as AccountService) getById(accountId int) (utils.Account, error) {
	account, err := as.accountDb.GetById(accountId)

	if err != nil {
		return utils.Account{}, errors.New("Failed to get account with id: " + strconv.Itoa(accountId))
	}

	return account, nil
}

type AccountService struct {
	accountDb utils.AccountDatabase
}

func validatePassword(password string, hashNsalt string) (bool, error) {
	hash, salt, ok := strings.Cut(hashNsalt, ":")

	if !ok {
		return false, errors.New("Failed to split has and salt")
	}

	newHash := strongishHash(password, salt)

	return hash == newHash, nil
}

// Don't try this at home kids. Use a cloud service for authentication stuff.
func strongishHash(password string, salt string) string {
	return string(pbkdf2.Key([]byte(password), []byte(salt), 1000, 64, sha512.New))
}

func hashPassword(password string) string {
	randomBytes := make([]byte, 128)
	_, err := rand.Read(randomBytes)

	if err != nil {
		panic("Failed to generate random string: " + err.Error())
	}

	salt := string(randomBytes)
	hash := strongishHash(password, string(salt))

	return hash + ":" + salt
}
