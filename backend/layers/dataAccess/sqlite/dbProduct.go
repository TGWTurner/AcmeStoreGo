package sqlite

import (
	"bjssStoreGo/backend/utils"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Id                int      `gorm:"primaryKey"`
	QuantityRemaining int      `gorm:"not null"`
	CategoryId        int      `gorm:"index;not null"`
	Price             int      `gorm:"not null"`
	ShortDescription  string   `gorm:"not null"`
	LongDescription   string   `gorm:"not null"`
	Category          Category `gorm:"ForeignKey:CategoryId"`
}

func ConvertToDbProduct(product utils.Product) Product {
	return Product{
		Id:                product.Id,
		QuantityRemaining: product.QuantityRemaining,
		CategoryId:        product.CategoryId,
		Price:             product.Price,
		ShortDescription:  product.ShortDescription,
		LongDescription:   product.LongDescription,
	}
}

func ConvertToDbProducts(products []utils.Product) []Product {
	dbProduct := make([]Product, len(products))

	for i, product := range products {
		dbProduct[i] = ConvertToDbProduct(product)
	}

	return dbProduct
}

func ConvertFromDbProduct(product Product) utils.Product {
	return utils.Product{
		Id:                product.Id,
		QuantityRemaining: product.QuantityRemaining,
		CategoryId:        product.CategoryId,
		Price:             product.Price,
		ShortDescription:  product.ShortDescription,
		LongDescription:   product.LongDescription,
	}
}

func ConvertFromDbProducts(dbProducts []Product) []utils.Product {
	products := make([]utils.Product, len(dbProducts))

	for i, product := range dbProducts {
		products[i] = ConvertFromDbProduct(product)
	}

	return products
}
