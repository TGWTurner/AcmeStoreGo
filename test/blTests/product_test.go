package blTests

import (
	"backend/layers/businessLogic"
	"backend/layers/dataAccess/testData"
	"backend/utils"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestSearchProductsForDealsReturnsCurrentDeals(t *testing.T) {
	ps := SetUpProduct()
	defer ps.Close()

	query := map[string]string{
		"dealDate": utils.GetFormattedDate(),
	}

	productsWithDeals, err := ps.SearchProducts(query)

	AssertNil(t, err)

	expectedProducts := getProductsWithDeals()

	if !assertProductSlicesAreEqual(productsWithDeals, expectedProducts) {
		t.Errorf("Failed to provide correct deals")
	}
}

func TestSearchProductsForCategoryReturnsCorrectProducts(t *testing.T) {
	ps := SetUpProduct()
	defer ps.Close()

	categoryId := 1

	query := map[string]string{
		"category": strconv.Itoa(categoryId),
	}

	productsFromCategory, err := ps.SearchProducts(query)

	AssertNil(t, err)

	expectedProducts := getProductsInCategory(categoryId)

	if !assertProductSlicesAreEqual(productsFromCategory, expectedProducts) {
		t.Errorf("Failed to return correct products in category")
	}
}

func TestSearchProductsForTextReturnsCorrectProducts(t *testing.T) {
	ps := SetUpProduct()
	defer ps.Close()

	searchTerm := "dog"

	query := map[string]string{
		"search": searchTerm,
	}

	productsFromSearch, err := ps.SearchProducts(query)

	AssertNil(t, err)

	expectedProducts := getProductsWithText(searchTerm)

	if !assertProductSlicesAreEqual(productsFromSearch, expectedProducts) {
		t.Errorf("Failed to return correct products for search text")
	}
}

func TestSearchProductsWithoutQueryReturnsAllProducts(t *testing.T) {
	ps := SetUpProduct()
	defer ps.Close()

	query := map[string]string{}

	productsFromSearch, err := ps.SearchProducts(query)

	AssertNil(t, err)

	expectedProducts := testData.GetProductTestData().Products

	if !assertProductSlicesAreEqual(productsFromSearch, expectedProducts) {
		t.Errorf("Failed to return all products")
	}
}

func TestGetProductCategoriesReturnsCategories(t *testing.T) {
	ps := SetUpProduct()
	defer ps.Close()

	categories, err := ps.GetProductcategories()

	AssertNil(t, err)

	expectedCategories := testData.GetProductTestData().Categories

	if !reflect.DeepEqual(categories, expectedCategories) {
		t.Errorf("Failed to return correct categories")
	}
}

func TestCheckStockReturnsCorrectProductsWithoutEnoughStock(t *testing.T) {
	ps := SetUpProduct()
	defer ps.Close()

	productWithoutEnoughStock := getTestProductFromId(5)

	orderItems := []utils.OrderItem{
		{
			ProductId: 1,
			Quantity:  1,
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

	AssertNil(t, err)

	if !reflect.DeepEqual(notEnoughStock, expectedOrderItems) {
		t.Errorf("Failed to return expected product as without enough stock")
	}
}

func TestCheckStockReturnsNoProductsWhenAllHaveStock(t *testing.T) {
	ps := SetUpProduct()
	defer ps.Close()

	orderItems := []utils.OrderItem{
		{
			ProductId: 1,
			Quantity:  1,
		},
		{
			ProductId: 5,
			Quantity:  1,
		},
	}

	notEnoughStock, _, err := ps.CheckStock(orderItems)

	AssertNil(t, err)

	if len(notEnoughStock) != 0 {
		t.Errorf("Failed to return no products without enough stock")
	}
}

func TestCheckStockReturnsCorrectTotal(t *testing.T) {
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

	AssertNil(t, err)

	expectedTotal := calculateTotalFromOrderItems(orderItems)

	if total != expectedTotal {
		t.Errorf("Expected total: %d, got: %d", expectedTotal, total)
	}
}

func TestDecreaseStockReducesStockOfProvidedProducts(t *testing.T) {
	ps := SetUpProduct()
	defer ps.Close()

	orderItems := []utils.OrderItem{
		{
			ProductId: 1,
			Quantity:  1,
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

	AssertNil(t, err)

	for _, item := range orderItems {
		alteredProduct, err := getPsProductFromId(ps, item.ProductId)

		AssertNil(t, err)

		for _, product := range products {
			if product.Id != alteredProduct.Id {
				continue
			}

			if product.QuantityRemaining != alteredProduct.QuantityRemaining {
				t.Errorf("Failed to change stock level to correct value")
			}
		}
	}
}

func TestDecreaseStockReturnsErrorWhenProductDoesNotHaveEnoughStock(t *testing.T) {
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
		t.Errorf("Failed to return error when products did not have enough stock")
	}
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
