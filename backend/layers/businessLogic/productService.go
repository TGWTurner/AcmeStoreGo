package businessLogic

import "bjssStoreGo/backend/utils"

func NewProductService(productDatabase utils.ProductDatabase) *ProductService {
	return &ProductService{
		db: productDatabase,
	}
}

func (ps ProductService) searchProducts() {
	//TODO: Implement search products
}

func (ps ProductService) getProductcategories() {
	//TODO: Implement get product categories
}

func (ps ProductService) checkStock() {
	//TODO: Implement check stock
}

func (ps ProductService) decreaseStock() {
	//TODO: Implement decrease stock
}

type ProductService struct {
	db utils.ProductDatabase
}
