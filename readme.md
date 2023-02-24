# BJSS Store 

This is a Go backend using a package called Gorilla Mux to help implement the APIs. 

# Getting started
Assumes you are running on 
* Windows and have go installed (https://go.dev/doc/install)
* Visual Studio Code

```
cd backend
code .
```
Things to try:
* Start the server `go run main.go`, poke around here: http://localhost:4002/api-docs
* Run the unit tests `go test ./...`
* Change the environment variable `DB_CONNECTION` to choose database options `sql`/`sql-mem`/``

** The integration tests also run the server. **

# Debugging

In VS Code:
* Set a breakpoint by clicking to the right of the line number of the line of code you want to the debugger to stop on. A red dot will appear. 
* Go to the run/debug tab. In the drop down select `Node.js...` then npm run script you want to debug. Try `Run Script: test`
If you choose to debug the server you'll have to make requests to it somehow.
* Assuming it reaches the line of code you set a breakpoint on, you can now hover over variables to see what they are set to and use the controls at the top to step to the next line or step into a function or run to the next breakpoint. 

# Databases
By default it runs with a simple in-memory data store. To use SQLite set the
`DB_CONNECTION` environment variable to `sql`. To use SQLite in memory set to `sql-mem`. To use a simple memory store do not set the variable.
For example for SQLite:
```
$env:DB_CONNECTION="sql"; go run test ./...
```
OR for SQLite Memory
```
$env:DB_CONNECTION="sql-mem"; go run test ./...
```
OR for Simple Memory Store
```
go run test ./...
```
At the moment it deletes the database at the start of every run to make life simpler and nudge us towards automated rather than manual testing.  You can examine the contents of the SQL database by using the VS Code extension SQLite.
 * Install the extension into your VS Code: https://marketplace.visualstudio.com/items?itemName=alexcvzz.vscode-sqlite
 * Right click on the created sqlite.db file, selecting `Open Database`
 * Inspect created tables
 * To run a query: 
    * ctrl + shift + p
    * SQLite: Quick Query
    * Select database file `sqlite.db`
    * Write query eg: `SELECT * FROM categories`

# Design Notes (WIP)
Useful types (From Backend/utils/utils.go and Backend/layers/api/apiUtils.go)
```go
type ShippingDetails struct {
	Email string, Name string, Address  string, Postcode string
}

type Account struct {
	Id int, PasswordHash string, ShippingDetails
}

type AccountDetails struct {
	Password string, utils.ShippingDetails
}

type AccountApiResponse struct {
	Id int, ShippingDetails
}

type Product struct {
	Id int, QuantityRemaining int, CategoryId int, Price int, ShortDescription string, LongDescription string
}

type ProductCategory struct {
	Id int, Name string
}

type ProductDeal struct {
	ProductId int, StartDate string, EndDate string
}

type Order struct {
	Id string, Total int, UpdatedDate string, CustomerId int, ShippingDetails ShippingDetails, Items []OrderItem
}

type OrderItem struct {
	ProductId int, Quantity  int
}

type OrderRequest struct {
	PaymentToken string, ShippingDetails utils.ShippingDetails, Items []utils.OrderItem
}

type OrderApiResponse struct {
	Id string, Total int, UpdatedDate string, ShippingDetails utils.ShippingDetails, Items []utils.OrderItem
}

type Basket struct {
	Total int, Items []OrderItem
}

type ApiErrorResponse struct {
	Error string, Msg string
}
```

API Model Summary. Run and access /api-docs/ or look at the /docs/docs.go file for more details
```
POST /account/sign-in {email: string, password: string} => AccountApiResponse, starts session
POST /account/sign-up ShippingDetails => AccountApiResponse, starts session
GET  /account ShippingDetails => AccountApiResponse, must be signed in
POST /account ShippingDetails => AccountApiResponse, must be signed in
GET  /product/catalogue?search|category [Product]
GET  /product/categories [ProductCategory]
GET  /product/deals [Product]
GET  /order/basket Basket
POST /order/basket Basket => Basket, starts session
GET  /order/history [Order], must be signed in
POST /order/checkout OrderRequest => OrderResponse, starts session
```

Desired user behaviour is
- Guest user can view products, deals categories, create/view basket, place order
- Guest user signing up should not empty basket 
- Customer (a signed in user) can see basket across devices.
- Guest or Customer can following link in email sent after order to see it without signing in

Other design notes
- Using Go 1.19.4 windows/amd64
- Uses Go package swaggo to produce the js/html page /api-docs/. This package uses the doc comments to generate the /docs folder and api spec.
