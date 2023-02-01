package businessLogic

import (
	"bjssStoreGo/backend/utils"
	"errors"
	"strconv"
)

func NewProductService(productDatabase utils.ProductDatabase) *ProductService {
	return &ProductService{
		productDb: productDatabase,
	}
}

func (ps ProductService) SearchProducts(query map[string]string) ([]utils.Product, error) {
	if len(query) > 0 {
		if val, ok := query["dealDate"]; ok {
			return ps.productDb.GetWithCurrentDeals(val)
		} else if val, ok := query["dealDate"]; ok {
			categoryId, err := strconv.Atoi(val)

			if err != nil {
				return []utils.Product{}, errors.New("Failed to convert categoryId to int")
			}

			return ps.productDb.GetByCategory(categoryId)
		} else if val, ok := query["dealDate"]; ok {
			return ps.productDb.GetByText(val)
		}
	}

	return ps.productDb.GetAll()
}

func (ps ProductService) GetProductcategories() ([]utils.ProductCategory, error) {
	return ps.productDb.GetCategories()
}

func (ps ProductService) CheckStock(productQuantities []utils.OrderItem) ([]utils.OrderItem, int, error) {
	total, err := ps.calculateTotalFromOrderItems(productQuantities)

	if err != nil {
		return []utils.OrderItem{}, 0, err
	}

	notEnoughStock, err := ps.calculateProductsLackingStockFromOrderItems(productQuantities)

	if err != nil {
		return []utils.OrderItem{}, 0, err
	}

	return notEnoughStock, total, err
}

func (ps ProductService) calculateProductsLackingStockFromOrderItems(orderItems []utils.OrderItem) ([]utils.OrderItem, error) {
	var notEnoughStock []utils.OrderItem

	for _, orderItem := range orderItems {
		product, err := ps.productDb.GetById(orderItem.ProductId)

		if err != nil {
			return []utils.OrderItem{}, err
		}

		if orderItem.Quantity > product.QuantityRemaining {
			notEnoughStock = append(
				notEnoughStock,
				utils.OrderItem{
					ProductId: orderItem.ProductId,
					Quantity:  product.QuantityRemaining,
				},
			)
		}
	}

	return notEnoughStock, nil
}

func (ps ProductService) calculateTotalFromOrderItems(orderItems []utils.OrderItem) (int, error) {
	var total int

	for _, orderItem := range orderItems {
		product, err := ps.productDb.GetById(orderItem.ProductId)

		if err != nil {
			return 0, err
		}

		total = total + (product.Price * orderItem.Quantity)
	}

	return total, nil
}

func (ps ProductService) DecreaseStock(productQuantities []utils.OrderItem) error {
	notEnoughStock, _, err := ps.CheckStock(productQuantities)

	if err != nil {
		return err
	}

	if len(notEnoughStock) > 0 {
		return errors.New("Trying to decrease stock below zero")
	}

	return ps.productDb.DecreaseStock(productQuantities)
}

type ProductService struct {
	productDb utils.ProductDatabase
}
