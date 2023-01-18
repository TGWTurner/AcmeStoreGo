package businessLogic

func NewAccountService(accountDb string) AccountService {
	return AccountService{
		accountDb: accountDb,
	}
}

func (as AccountService) signIn(email string, password string) {
	//TODO: Implement sign in
}

func (as AccountService) signUp(email string, password string, name string, address string, postcode string) {
	//TODO: Implement sign up
}

func (as AccountService) update(account string) {
	//TODO: Implement update
}

func (as AccountService) getById(accountId int) {
	//TODO: Implement get by id
}

type AccountService struct {
	accountDb string //TODO: UPDATE
}
