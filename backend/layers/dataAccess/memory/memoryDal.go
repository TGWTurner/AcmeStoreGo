package memory

func NewDatabase() Database {
	return Database{
		Account: NewAccountDatabase(),
		Product: NewProductDatabase(),
		Order:   NewOrderDatabase(),
	}
}

type Database struct {
	Account AccountDatabase
	Product ProductDatabase
	Order   OrderDatabase
}
