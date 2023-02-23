// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/account": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Gets the user's Account",
                "operationId": "GetAccount",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token",
                        "name": "Cookie",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The user's account",
                        "schema": {
                            "$ref": "#/definitions/utils.AccountApiResponse"
                        }
                    },
                    "401": {
                        "description": "User is not signed in",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Updates the users account",
                "operationId": "PostAccount",
                "parameters": [
                    {
                        "description": "account information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.AccountDetails"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The user's account",
                        "schema": {
                            "$ref": "#/definitions/utils.AccountApiResponse"
                        }
                    },
                    "401": {
                        "description": "User is not signed in",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/account/sign-in": {
            "post": {
                "description": "Signs in, deletes any existing session, creates a new one for this user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Signs in",
                "operationId": "PostSignIn",
                "parameters": [
                    {
                        "description": "account information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UserDetails"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The user's account",
                        "schema": {
                            "$ref": "#/definitions/utils.AccountApiResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/account/sign-up": {
            "post": {
                "description": "Signs up, deletes any existing session, creates a new one for this user. Will give an error if the user already exists.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Registers a new user",
                "operationId": "PostSignUp",
                "parameters": [
                    {
                        "description": "account information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.AccountDetails"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The user's account",
                        "schema": {
                            "$ref": "#/definitions/utils.AccountApiResponse"
                        }
                    },
                    "400": {
                        "description": "Malformed request or account already exists",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/order/basket": {
            "get": {
                "description": "The same session cookie that created the basket is needed",
                "produces": [
                    "application/json"
                ],
                "summary": "Gets the user's Basket",
                "operationId": "GetBasket",
                "responses": {
                    "200": {
                        "description": "A Basket",
                        "schema": {
                            "$ref": "#/definitions/utils.Basket"
                        }
                    }
                }
            },
            "post": {
                "description": "Sets a session cookie which is needed to later get the basket",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Creates or updates the user's Basket",
                "operationId": "PostBasket",
                "parameters": [
                    {
                        "description": "A Basket",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.Basket"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A Basket",
                        "schema": {
                            "$ref": "#/definitions/utils.Basket"
                        }
                    }
                }
            }
        },
        "/api/order/checkout": {
            "post": {
                "description": "Checks the stock levels and paymentToken. If Ok creates a new order. If not gives an error and the products there is not enough stock for. Sets a session cookie which can be used later to tie this order to a signed in user. Does not require a signed in user so guests can check out",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Creates a new order",
                "operationId": "PostCheckout",
                "parameters": [
                    {
                        "description": "An order",
                        "name": "order",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.Order"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A Basket",
                        "schema": {
                            "$ref": "#/definitions/utils.Order"
                        }
                    },
                    "400": {
                        "description": "An error. If the request was well formed this will be payment or stock level error. If stock level error, the quantityRemaining is returned for products with not enough stock.",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/order/history": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Gets the user's Order history",
                "operationId": "GetHistory",
                "responses": {
                    "200": {
                        "description": "A Basket",
                        "schema": {
                            "$ref": "#/definitions/utils.Basket"
                        }
                    }
                }
            }
        },
        "/api/order/{token}": {
            "get": {
                "description": "Does not require a signed in user so that we can implement getting an order via a link in an email, etc.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Fetches an order given a token",
                "operationId": "GetOrder",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token",
                        "name": "Cookie",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A Basket",
                        "schema": {
                            "$ref": "#/definitions/utils.Order"
                        }
                    },
                    "404": {
                        "description": "No such order",
                        "schema": {
                            "$ref": "#/definitions/utils.ApiErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/product/catalogue": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Query or get all Products",
                "operationId": "Search",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Text to search for",
                        "name": "search",
                        "in": "path"
                    },
                    {
                        "type": "integer",
                        "description": "A Category Id to filter Products on",
                        "name": "caregory",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "A JSON array of Products",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/utils.Product"
                            }
                        }
                    }
                }
            }
        },
        "/api/product/categories": {
            "get": {
                "description": "Signs up, deletes any existing session, creates a new one for this user. Will give an error if the user already exists.",
                "produces": [
                    "application/json"
                ],
                "summary": "Get a list of product Categories",
                "operationId": "Categories",
                "responses": {
                    "200": {
                        "description": "A JSON array of Categories",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/utils.ProductCategory"
                            }
                        }
                    }
                }
            }
        },
        "/api/product/deals": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get deals that are valid for today",
                "operationId": "Deals",
                "responses": {
                    "200": {
                        "description": "A JSON array of Products",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/utils.Product"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.UserDetails": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "utils.AccountApiResponse": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "postcode": {
                    "type": "string"
                }
            }
        },
        "utils.AccountDetails": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "postcode": {
                    "type": "string"
                }
            }
        },
        "utils.ApiErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "utils.Basket": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/utils.OrderItem"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "utils.Order": {
            "type": "object",
            "properties": {
                "customerId": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/utils.OrderItem"
                    }
                },
                "shippingDetails": {
                    "$ref": "#/definitions/utils.ShippingDetails"
                },
                "total": {
                    "type": "integer"
                },
                "updatedDate": {
                    "type": "string"
                }
            }
        },
        "utils.OrderItem": {
            "type": "object",
            "properties": {
                "productId": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "utils.Product": {
            "type": "object",
            "properties": {
                "categoryId": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "longDescription": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "quantityRemaining": {
                    "type": "integer"
                },
                "shortDescription": {
                    "type": "string"
                }
            }
        },
        "utils.ProductCategory": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "utils.ShippingDetails": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "postcode": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:4001",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "BJSS Store",
	Description:      "Simple store for teaching and learning",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
