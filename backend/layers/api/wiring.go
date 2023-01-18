/*
TODO:
 - API wiring method
 - Look at the one in the academy go store
 - Need to figure out some form of session storage
 - Requires routes for:
  - /api/account/sign-in: POST
  - /api/account/sign-up: POST
  - /api/account: GET
  - /api/account: POST

  - /api/product/catalogue: GET
  - /api/product/categories: GET
  - /api/product/deals: GET

  - /api/order/basket: GET
  - /api/order/basket: POST
  - /api/order/checkout: POST
  - /api/order/history: GET

  - /api/order/:id GET
  and a 404

*/

package api

type wiring struct{}
