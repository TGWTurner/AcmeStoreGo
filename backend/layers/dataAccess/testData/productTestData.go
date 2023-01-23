package testData

import (
	"bjssStoreGo/backend/utils"
)

type Data struct {
	Products   []utils.Product
	Categories []utils.ProductCategory
	Deals      []utils.ProductDeal
}

func (ptd ProductTestData) GetTestData() Data {

	return Data{
		Products: []utils.Product{
			{
				Id:                6,
				CategoryId:        1,
				Price:             100,
				QuantityRemaining: 2,
				ShortDescription:  "Dog",
				LongDescription:   "The dog (Canis familiaris when considered a distinct species or Canis lupus familiaris when considered a subspecies of the wolf) is a domesticated carnivore of the family Canidae...",
			},
			{
				Id:                1,
				CategoryId:        1,
				Price:             1000,
				QuantityRemaining: 1000,
				ShortDescription:  "Giraffe",
				LongDescription:   "The giraffe (Giraffa) is an African artiodactyl mammal, the tallest living terrestrial animal and the largest ruminant. It is traditionally considered to be one species, Giraffa...",
			},
			{
				Id:                2,
				CategoryId:        1,
				Price:             90,
				QuantityRemaining: 1000,
				ShortDescription:  "Koala",
				LongDescription:   "The koala or, inaccurately, koala bear[a] (Phascolarctos cinereus) is an arboreal herbivorous marsupial native to Australia. It is the only extant representative of the family Phascolarctidae...",
			},
			{
				Id:                3,
				CategoryId:        2,
				Price:             1,
				QuantityRemaining: 2,
				ShortDescription:  "Brazil Nut",
				LongDescription:   "The Brazil nut (Bertholletia excelsa) is a South American tree in the family Lecythidaceae, and it is also the name of the trees commercially harvested edible seeds. It is one of the largest...",
			},
			{
				Id:                4,
				CategoryId:        2,
				Price:             2,
				QuantityRemaining: 2,
				ShortDescription:  "Apricot",
				LongDescription:   "An apricot is a fruit, or the tree that bears the fruit, of several species in the genus Prunus (stone fruits).",
			},
			{
				Id:                5,
				CategoryId:        2,
				Price:             3,
				QuantityRemaining: 2,
				ShortDescription:  "Orange",
				LongDescription:   "The orange is the fruit of various citrus species in the family Rutaceae (see list of plants known as orange); it primarily refers to Citrus × sinensis, which is also called sweet orange...",
			},
		},
		Categories: []utils.ProductCategory{
			{
				Id:   1,
				Name: "Animals",
			},
			{
				Id:   2,
				Name: "Fruits",
			},
		},
		Deals: []utils.ProductDeal{
			{
				ProductId: 1,
				StartDate: "2021-02-13",
				EndDate:   "2030-02-13",
			},
			{
				ProductId: 2,
				StartDate: "2021-02-13",
				EndDate:   "2030-02-13",
			},
			{
				ProductId: 4,
				StartDate: "1970-01-13",
				EndDate:   "1970-01-13",
			},
		},
	}
}

type ProductTestData struct{}
