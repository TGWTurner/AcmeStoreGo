package blTests

import (
	"bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/layers/dataAccess"
	"testing"
)

func SetUpAccount() businessLogic.AccountService {
	db := dataAccess.InitiateConnection()
	accountService := businessLogic.NewAccountService(db.Account)

	return *accountService
}

func SetUpOrder() businessLogic.OrderService {
	db := dataAccess.InitiateConnection()
	ps := SetUpProduct()
	orderService := businessLogic.NewOrderService(db.Order, ps)

	return *orderService
}

func SetUpProduct() businessLogic.ProductService {
	db := dataAccess.InitiateConnection()
	productService := businessLogic.NewProductService(db.Product)

	return *productService
}

func AssertNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected error: nil, got: %s", err.Error())
	}
}
