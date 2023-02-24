package sqlite

import (
	"backend/utils"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Id   int    `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func ConvertToDbCategory(category utils.ProductCategory) Category {
	return Category{
		Id:   category.Id,
		Name: category.Name,
	}
}

func ConvertToDbCategories(categories []utils.ProductCategory) []Category {
	dbCategories := make([]Category, len(categories))

	for i, category := range categories {
		dbCategories[i] = ConvertToDbCategory(category)
	}

	return dbCategories
}

func ConvertFromDbCategory(category Category) utils.ProductCategory {
	return utils.ProductCategory{
		Id:   category.Id,
		Name: category.Name,
	}
}

func ConvertFromDbCategories(dbCategories []Category) []utils.ProductCategory {
	categories := make([]utils.ProductCategory, len(dbCategories))

	for i, category := range dbCategories {
		categories[i] = ConvertFromDbCategory(category)
	}

	return categories
}
