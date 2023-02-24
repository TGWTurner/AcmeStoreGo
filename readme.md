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
`DB_CONNECTION` environment variable to `sql`, for example:
```
DB_CONNECTION=sql npm run test-all
```
or 
```
DB_CONNECTION=sql npm start
```
At the moment it deletes the database at the start of every run to make life simpler and nudge us towards automated rather than manual testing.  You can examine the contents of the SQL database by loading it into a tool such as this: https://chrome.google.com/webstore/detail/sqlite-manager/njognipnngillknkhikjecpnbkefclfe or on the command line:
* Install `sudo apt-get install sqlite3`
* Run a query `sqlite3 sqlite.db "select * from categories"`

You can open an interactive session with `sqlite3 sqlite.db`, but *beware* its a static view of the database at the point in time you started the session. You won't see any updates after. Example session
```
sqlite> .tables
accounts    categories  deals       productssele
sqlite> select * from categories;
--- data ---
sqlite> .quit
```

# Design Notes (WIP)
Useful types
```typescript
type ShippingDetails = {
    email: string, name: string, address: string, postcode: string
}
type Account = { id: string, passwordHash: string } & ShippingDetails
type AccountApiResponse = {id: string} & ShippingDetails

type Product = {
    id: number, 
    quantityRemaining: number
    categoryId: number, 
    price: number,
    shortDescription: string, 
    longDescription: string
}
type ProductCategory = { id: number, name: string }
type ProductDeal = { productId: number, startDate: Date, endDate: Date }

type OrderItem = { productId: number, quantity: number }
type Order = {
    id: string, 
    total: number, 
    updatedDate: Date
    customerId?: number, // not present on guest checkout
    shippingDetails: ShippingDetails, 
    items: [OrderItem],
}
type OrderApiRequest = {
    paymentToken: string, 
    shippingDetails: ShippingDetails, 
    items: [OrderItem]
}
type OrderApiResponse = Order | {
    error: string, 
    items: [{productId: number, quantityRemaining: number}]
}
type Basket = {total: number, items: [OrderItem]}

type Session = {basket?: Basket, customerId: string}
```

API Model Summary.  See the OpenAPI spec for details
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
- Using Node 12+ for ES6 modules (`import`).
- Not Typescript for simplicity of setup. Also it's a more complex language to learn. However, you can see by the Typescript model above that I need some type definitions to help me reason about the code, so I'm on the fence about if JS is actually any easier.
- Express rather than Azure Functions. Mainly because it is easier to Google for Express solutions, but also because sessions/authorisation would need a re-think.
- No 'code first' OpenAPI spec despite the tediousness of OpenAPI to write. Mainly for simplicity, but also because the spec might be useful a useful artefact to test other implementations against. 
