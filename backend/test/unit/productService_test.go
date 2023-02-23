package unit

import (
	bl "bjssStoreGo/backend/layers/businessLogic"
	da "bjssStoreGo/backend/layers/dataAccess"
	"bjssStoreGo/backend/layers/dataAccess/testData"
	"bjssStoreGo/backend/test"
	"bjssStoreGo/backend/utils"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func setUpProduct(t *testing.T) bl.ProductService {
	//"sql" or "sql-mem" or ""
	t.Setenv("DB_CONNECTION", "sql-mem")

	db := da.InitiateConnection()
	return *bl.NewProductService(db.Product)
}

func TestGetsAll(t *testing.T) {
	ps := setUpProduct(t)
	defer ps.Close()

	products, err := ps.SearchProducts(map[string]string{})

	test.AssertNil(t, err)

	test.AssertProductSetsMatch(t, products, testData.GetProductTestData().Products)
}

func TestFindsByCategory(t *testing.T) {
	ps := setUpProduct(t)
	defer ps.Close()

	query := map[string]string{"category": "1"}
	products, err := ps.SearchProducts(query)

	test.AssertNil(t, err)

	expected := getProductsInCategory(1)

	test.AssertProductSetsMatch(t, products, expected)
}

func TestFindsByText(t *testing.T) {
	ps := setUpProduct(t)
	defer ps.Close()

	query := map[string]string{"search": "Apricot"}
	products, err := ps.SearchProducts(query)

	test.AssertNil(t, err)

	expected := getProductsWithText("Apricot")

	test.AssertProductSetsMatch(t, products, expected)

	query = map[string]string{"search": "fruit"}
	products, err = ps.SearchProducts(query)

	test.AssertNil(t, err)

	expected = getProductsWithText("fruit")

	test.AssertProductSetsMatch(t, products, expected)
}

func TestGetsInDateDeals(t *testing.T) {
	ps := setUpProduct(t)
	defer ps.Close()

	query := map[string]string{"dealDate": utils.GetFormattedDate()}
	products, err := ps.SearchProducts(query)

	test.AssertNil(t, err)

	expected := getProductsWithDeals(utils.GetFormattedDate())

	test.AssertProductSetsMatch(t, products, expected)
}

func TestGetsNoDealsIfNoneInDate(t *testing.T) {
	ps := setUpProduct(t)
	defer ps.Close()

	query := map[string]string{"dealDate": "2000-02-21"}
	products, err := ps.SearchProducts(query)

	test.AssertNil(t, err)

	expected := getProductsWithDeals("2000-02-21")

	test.AssertProductSetsMatch(t, products, expected)
}

func TestGetsCategories(t *testing.T) {
	ps := setUpProduct(t)
	defer ps.Close()

	categories, err := ps.GetProductcategories()

	test.AssertNil(t, err)

	expected := testData.GetProductTestData().Categories

	test.AssertCategorySetsMatch(t, categories, expected)
}

func TestCalcsTotalsAndStockShortages(t *testing.T) {
	ps := setUpProduct(t)
	defer ps.Close()

	orderItems := []utils.OrderItem{
		{ProductId: 1, Quantity: 1}, //One dog at 100
		{ProductId: 3, Quantity: 1}, //One koala at 90
		{ProductId: 5, Quantity: 4}, //Four apricots at 2
	}

	notEnoughStock, total, err := ps.CheckStock(orderItems)

	test.AssertNil(t, err)

	expectedTotal := 198

	if total != expectedTotal {
		t.Errorf("Total incorrect, expected: %d, actual: %d", expectedTotal, total)
	}

	expectedNotEnoughStock := []utils.OrderItem{
		{ProductId: 4, Quantity: 2}, //Only two apricots remain
	}

	if reflect.DeepEqual(notEnoughStock, expectedNotEnoughStock) {
		t.Errorf("Not enough stock products incorrect")
	}
}

func TestDecreasesStock(t *testing.T) {
	ps := setUpProduct(t)
	defer ps.Close()

	orderItems := []utils.OrderItem{
		{ProductId: 1, Quantity: 1},
		{ProductId: 3, Quantity: 2},
	}

	ps.DecreaseStock(orderItems)

	if product, err := getPsProductFromId(ps, 1); err != nil || product.QuantityRemaining != 1 {
		t.Errorf("Expected quantity of %d, got %d for product with id %d", 1, product.QuantityRemaining, product.Id)
	}

	if product, err := getPsProductFromId(ps, 3); err != nil || product.QuantityRemaining != 998 {
		t.Errorf("Expected quantity of %d, got %d for product with id %d", 998, product.QuantityRemaining, product.Id)
	}

	if product, err := getPsProductFromId(ps, 2); err != nil || product.QuantityRemaining != 1000 {
		t.Errorf("Expected quantity of %d, got %d for product with id %d", 1000, product.QuantityRemaining, product.Id)
	}
}

func TestDoesNotDecreaseStockBelow0(t *testing.T) {
	ps := setUpProduct(t)
	defer ps.Close()

	orderItems := []utils.OrderItem{
		{ProductId: 2, Quantity: 10000},
	}

	err := ps.DecreaseStock(orderItems)

	if err == nil {
		t.Errorf("Expected error got none")
	}
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

func getDeals(date string) []utils.ProductDeal {
	allDeals := testData.GetProductTestData().Deals
	currentDeals := []utils.ProductDeal{}

	for _, deal := range allDeals {
		if deal.EndDate > date && deal.StartDate < date {
			currentDeals = append(currentDeals, deal)
		}
	}

	return currentDeals
}

func getProductsWithDeals(date string) []utils.Product {
	currentDeals := getDeals(date)
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

func getPsProductFromId(ps bl.ProductService, productId int) (utils.Product, error) {
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
