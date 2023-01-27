package memory

import "bjssStoreGo/backend/utils"

func NewDatabase() utils.Database {
	return utils.Database{
		Account: NewAccountDatabase(),
		Product: NewProductDatabase(),
		Order:   NewOrderDatabase(),
	}
}
