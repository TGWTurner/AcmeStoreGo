package sqlite

import (
	"backend/utils"

	"gorm.io/gorm"
)

type Deal struct {
	gorm.Model
	ProductId int     `gorm:"not null"`
	StartDate string  `gorm:"not null"`
	EndDate   string  `gorm:"not null"`
	Product   Product `gorm:"ForeignKey:ProductId"`
}

func ConvertToDbDeal(deal utils.ProductDeal) Deal {
	return Deal{
		ProductId: deal.ProductId,
		StartDate: deal.StartDate,
		EndDate:   deal.EndDate,
	}
}

func ConvertToDbDeals(deals []utils.ProductDeal) []Deal {
	dbDeals := make([]Deal, len(deals))

	for i, deal := range deals {
		dbDeals[i] = ConvertToDbDeal(deal)
	}

	return dbDeals
}

func ConvertFromDbDeal(deal Deal) utils.ProductDeal {
	return utils.ProductDeal{
		ProductId: deal.ProductId,
		StartDate: deal.StartDate,
		EndDate:   deal.EndDate,
	}
}

func ConvertFromDbDeals(dbDeals []Deal) []utils.ProductDeal {
	deals := make([]utils.ProductDeal, len(dbDeals))

	for i, deal := range dbDeals {
		deals[i] = ConvertFromDbDeal(deal)
	}

	return deals
}
