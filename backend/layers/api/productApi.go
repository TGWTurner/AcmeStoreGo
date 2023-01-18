package api

import (
	"bjssStoreGo/backend/layers/businessLogic"
)

func (p ProductApi) search(req string, res string) {
	//TODO: Implement search
	/*
			const products = await productService.searchProducts(req.query)
		        res.json(products)
	*/

	businessLogic.ProductService.SearchProducts()
}

func (p ProductApi) categories(req string, res string) {
	//TODO: Implement cataegories
	/*
			const categories = await productService.getProductCategories()
		        res.json(categories)
	*/
}

func (p ProductApi) deals(_, res string) {
	//TODO: Implement deals
	/*
		const products = await productService.searchProducts({ dealDate: new Date() })
			res.json(products)
	*/
}

type ProductApi struct{}
