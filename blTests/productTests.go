package blTests

import (
	"bjssStoreGo/backend/layers/businessLogic"
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/utils"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func TestSearchProductsForDealsReturnsCurrentDeals() (bool, string, string) {
	ps := SetUpProduct()
	defer ps.Close()

	query := map[string]string{
		"dealDate": utils.GetFormattedDate(),
	}

	productsWithDeals, err := ps.SearchProducts(query)

	if err != nil {
		return false, "testSearchProductsForDealsReturnsCurrentDeals", "Failed to get deals"
	}

	expectedProducts := getProductsWithDeals()

	if !assertProductSlicesAreEqual(productsWithDeals, expectedProducts) {
		return false, "testSearchProductsForDealsReturnsCurrentDeals", "Failed to provide correct deals"
	}

	return true, "testSearchProductsForDealsReturnsCurrentDeals", "Successfully provided correct deals"
}

func TestSearchProductsForCategoryReturnsCorrectProducts() (bool, string, string) {
	ps := SetUpProduct()
	defer ps.Close()

	categoryId := 1

	query := map[string]string{
		"category": strconv.Itoa(categoryId),
	}

	productsFromCategory, err := ps.SearchProducts(query)

	if err != nil {
		return false, "testSearchProductsForCategoryReturnsCorrectProducts", "Failed to search products"
	}

	expectedProducts := getProductsInCategory(categoryId)

	if !assertProductSlicesAreEqual(productsFromCategory, expectedProducts) {
		return false, "testSearchProductsForCategoryReturnsCorrectProducts", "Failed to return correct products in category"
	}

	return true, "testSearchProductsForCategoryReturnsCorrectProducts", "Successfully returned products in category"
}

func TestSearchProductsForTextReturnsCorrectProducts() (bool, string, string) {
	ps := SetUpProduct()
	defer ps.Close()

	searchTerm := "dog"

	query := map[string]string{
		"search": searchTerm,
	}

	productsFromSearch, err := ps.SearchProducts(query)

	if err != nil {
		return false, "testSearchProductsForTextReturnsCorrectProducts", "Failed to get products from search term"
	}

	expectedProducts := getProductsWithText(searchTerm)

	if !assertProductSlicesAreEqual(productsFromSearch, expectedProducts) {
		return false, "testSearchProductsForTextReturnsCorrectProducts", "Failed to return correct products for search text"
	}

	return true, "testSearchProductsForTextReturnsCorrectProducts", "Successfully returned correct products for search term"
}

func TestSearchProductsWithoutQueryReturnsAllProducts() (bool, string, string) {
	ps := SetUpProduct()
	defer ps.Close()

	query := map[string]string{}

	productsFromSearch, err := ps.SearchProducts(query)

	if err != nil {
		return false, "testSearchProductsWithoutQueryReturnsAllProducts", "Failed to search for products"
	}

	expectedProducts := testData.GetProductTestData().Products

	if !assertProductSlicesAreEqual(productsFromSearch, expectedProducts) {
		return false, "testSearchProductsWithoutQueryReturnsAllProducts", "Failed to return all products"
	}

	return true, "testSearchProductsWithoutQueryReturnsAllProducts", "Successfully returned all products "
}

func TestGetProductCategoriesReturnsCategories() (bool, string, string) {
	ps := SetUpProduct()
	defer ps.Close()

	categories, err := ps.GetProductcategories()

	if err != nil {
		return false, "testGetProductCategoriesReturnsCategories", "Failed to get product categories"
	}

	expectedCategories := testData.GetProductTestData().Categories

	if !reflect.DeepEqual(categories, expectedCategories) {
		return false, "testGetProductCategoriesReturnsCategories", "Failed to return correct categories"
	}

	return true, "testGetProductCategoriesReturnsCategories", "Successfully returned expected categories"
}

func TestCheckStockReturnsCorrectProductsWithoutEnoughStock() (bool, string, string) {
	ps := SetUpProduct()
	defer ps.Close()

	productWithoutEnoughStock := getTestProductFromId(5)

	orderItems := []utils.OrderItem{
		{
			ProductId: 1,
			Quantity:  5,
		},
		{
			ProductId: productWithoutEnoughStock.Id,
			Quantity:  productWithoutEnoughStock.QuantityRemaining + 1,
		},
	}

	expectedOrderItems := []utils.OrderItem{
		{
			ProductId: productWithoutEnoughStock.Id,
			Quantity:  productWithoutEnoughStock.QuantityRemaining,
		},
	}

	notEnoughStock, _, err := ps.CheckStock(orderItems)

	if err != nil {
		return false, "testCheckStockReturnsCorrectProductsWithoutEnoughStock", "Failed to check stock of products"
	}

	if !reflect.DeepEqual(notEnoughStock, expectedOrderItems) {
		return false, "testCheckStockReturnsCorrectProductsWithoutEnoughStock", "Failed to return expected product as without enough stock"
	}

	return true, "testCheckStockReturnsCorrectProductsWithoutEnoughStock", "Successfully returned correct product as without enough stock"
}

func TestCheckStockReturnsNoProductsWhenAllHaveStock() (bool, string, string) {
	ps := SetUpProduct()
	defer ps.Close()

	orderItems := []utils.OrderItem{
		{
			ProductId: 1,
			Quantity:  5,
		},
		{
			ProductId: 5,
			Quantity:  1,
		},
	}

	notEnoughStock, _, err := ps.CheckStock(orderItems)

	if err != nil {
		return false, "testCheckStockReturnsNoProductsWhenAllHaveStock", "Failed to check stock of products"
	}

	if len(notEnoughStock) != 0 {
		return false, "testCheckStockReturnsNoProductsWhenAllHaveStock", "Failed to return no products without enough stock"
	}

	return true, "testCheckStockReturnsNoProductsWhenAllHaveStock", "Successfully returned no products without enough stock"
}

func TestCheckStockReturnsCorrectTotal() (bool, string, string) {
	ps := SetUpProduct()
	defer ps.Close()

	orderItems := []utils.OrderItem{
		{
			ProductId: 1,
			Quantity:  5,
		},
		{
			ProductId: 5,
			Quantity:  1,
		},
	}

	_, total, err := ps.CheckStock(orderItems)

	if err != nil {
		return false, "testCheckStockReturnsCorrectTotal", "Failed to check stock of products"
	}

	expectedTotal := calculateTotalFromOrderItems(orderItems)

	if total != expectedTotal {
		return false, "testCheckStockReturnsCorrectTotal", "Failed to return no products without enough stock"
	}

	return true, "testCheckStockReturnsCorrectTotal", "Successfully returned no products without enough stock"
}

func TestDecreaseStockReducesStockOfProvidedProducts() (bool, string, string) {
	ps := SetUpProduct()
	defer ps.Close()

	orderItems := []utils.OrderItem{
		{
			ProductId: 1,
			Quantity:  5,
		},
		{
			ProductId: 5,
			Quantity:  1,
		},
	}

	products := []utils.Product{}
	for _, item := range orderItems {
		product := getTestProductFromId(item.ProductId)
		product.QuantityRemaining -= item.Quantity

		products = append(products, product)
	}

	err := ps.DecreaseStock(orderItems)

	if err != nil {
		return false, "testDecreaseStockReducesStockOfProvidedProducts", "Failed to decrease stock of products"
	}

	for _, item := range orderItems {
		alteredProduct, err := getPsProductFromId(ps, item.ProductId)

		if err != nil {
			return false, "testDecreaseStockReducesStockOfProvidedProducts", "Failed to get product from db"
		}

		for _, product := range products {
			if product.Id != alteredProduct.Id {
				continue
			}

			if product.QuantityRemaining != alteredProduct.QuantityRemaining {
				return false, "testDecreaseStockReducesStockOfProvidedProducts", "Failed to change stock level to correct value"
			}
		}
	}

	return true, "testDecreaseStockReducesStockOfProvidedProducts", "Successfully altered correct products remaining quantity"
}

func TestDecreaseStockReturnsErrorWhenProductDoesNotHaveEnoughStock() (bool, string, string) {
	ps := SetUpProduct()
	defer ps.Close()

	orderItems := []utils.OrderItem{
		{
			ProductId: 1,
			Quantity:  5,
		},
		{
			ProductId: 5,
			Quantity:  10,
		},
	}

	err := ps.DecreaseStock(orderItems)

	if err == nil {
		return false, "testDecreaseStockReturnsErrorWhenProductDoesNotHaveEnoughStock", "Failed to return error when products did not have enough stock"
	}

	return true, "testDecreaseStockReturnsErrorWhenProductDoesNotHaveEnoughStock", "Successfully returned error when attempting to reduce stock for products without sufficient stock"
}

func getPsProductFromId(ps businessLogic.ProductService, productId int) (utils.Product, error) {
	allProducts, err := ps.SearchProducts(map[string]string{})

	if err != nil {
		return utils.Product{}, err
	}

	for _, product := range allProducts {
		if product.Id == productId {
			return product, nil
		}
	}

	return utils.Product{}, errors.New("Product with id: " + strconv.Itoa(productId) + " not found")
}

func calculateTotalFromOrderItems(orderItems []utils.OrderItem) int {
	total := 0
	for _, item := range orderItems {
		product := getTestProductFromId(item.ProductId)

		total += product.Price * item.Quantity
	}

	return total
}

func getTestProductFromId(productId int) utils.Product {
	allProducts := testData.GetProductTestData().Products

	for _, product := range allProducts {
		if product.Id == productId {
			return product
		}
	}

	return utils.Product{}
}

func getProductsWithText(searchTerm string) []utils.Product {
	allProducts := testData.GetProductTestData().Products
	productsWithText := []utils.Product{}

	for _, product := range allProducts {
		if strings.Contains(product.ShortDescription, searchTerm) ||
			strings.Contains(product.LongDescription, searchTerm) {
			productsWithText = append(productsWithText, product)
		}
	}

	return productsWithText
}

func getProductsInCategory(categoryId int) []utils.Product {
	allProducts := testData.GetProductTestData().Products
	productsInCategory := []utils.Product{}

	for _, product := range allProducts {
		if product.CategoryId == categoryId {
			productsInCategory = append(productsInCategory, product)
		}
	}

	return productsInCategory
}

func getCurrentDeals() []utils.ProductDeal {
	currentDate := utils.GetFormattedDate()
	allDeals := testData.GetProductTestData().Deals
	currentDeals := []utils.ProductDeal{}

	for _, deal := range allDeals {
		if deal.EndDate > currentDate && deal.StartDate < currentDate {
			currentDeals = append(currentDeals, deal)
		}
	}

	return currentDeals
}

func getProductsWithDeals() []utils.Product {
	currentDeals := getCurrentDeals()
	allProducts := testData.GetProductTestData().Products
	currentDealProducts := []utils.Product{}

	for _, product := range allProducts {
		for _, deal := range currentDeals {
			if deal.ProductId == product.Id {
				currentDealProducts = append(currentDealProducts, product)
			}
		}
	}

	return currentDealProducts
}

func assertProductSlicesAreEqual(sliceA []utils.Product, sliceB []utils.Product) bool {
	if len(sliceA) != len(sliceB) {
		return false
	}

	for _, productA := range sliceA {
		matchFound := false

		for _, productB := range sliceB {
			if reflect.DeepEqual(productA, productB) {
				matchFound = true
				break
			}
		}

		if !matchFound {
			return false
		}
	}

	return true
}
