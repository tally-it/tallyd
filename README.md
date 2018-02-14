# Tally-It
A tool for hack spaces, clubs, groups and companies to manage a digital tally sheet.

## API Documentation

### User

| Method      | URL                         | Description                     |
| --------    | --------                    | --------                        |
| GET         | /v1/user                    | User index                      |
| POST        | /v1/user                    | Sign up                         |
| POST        | /v1/login                   | Login                           |
| GET         | /v1/user/:id                | User detail                     |
| POST        | /v1/user/:id/transaction    | Change balance +-               |
| PUT         | /v1/user/:id                | TODO: Change user data          |
| DELETE      | /v1/user/:id                | TODO: Delete user               |
| POST        | /v1/user/:id/auth/:authtype | TODO: Add authentication        |
| PUT         | /v1/user/:id/auth/:authtype | TODO: Change authentication     |
| DELETE      | /v1/user/:id/auth/:authtype | TODO: Delete authentication     |

#### User Index
Returns an index of all users.

- **Method**
    `GET`
- **URL**
    `/v1/user`
- **Request Body**
    `-`
- **Return Body**
    ```json
    [
        {
        "userID": 1,
        "name": "marove",
        "email": "blub@blub.de",
        "active": false,
        "isAdmin": false,
        "balance": 0
        },
        {
        "userID": 2,
        "name": "test",
        "email": "test@binary-kitchen.de",
        "active": false,
        "isAdmin": false,
        "balance": 0
        }
    ]
    ```
    
#### Sign Up
Adds user to database.

- **Method**
    `POST`
- **URL**
    `/v1/user`
- **Request Body**
    ```json
    {
        "name": "username",
        "email": "blub@blub.com",
        "password": "pa$$word"
    }
    ```
- **Return Body**
    ```json
    {
    "userID": 8
    }
    ```

#### Login
Login and create JWT.

- **Method**
    `POST`
    
- **URL**
    `/v1/login`
- **Request Body**
    ```json
    {
        "Name": "username",
        "Password": "pa$$word"
    }
    ```
- **Return Body**
   JWT-Token. Example: `"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOiIyMDE4LTAxLTI5VDExOjE0OjE5LjIzNjAwMjE2KzAxOjAwIiwiSXNBY3RpdmUiOmZhbHNlLCJJc0FkbWluIjpmYWxzZSwiTmFtZSI6InRlc3R1c2VyIiwidXNlcklEIjo2fQ.Vm2NKt8a5KXKQhEaeb1wBQbPplAXlrkhZ05ZgaKHIAY"`

#### User Detail
Returns user detail based on user id.

- **Method**
    `GET`
- **URL**
    `/v1/user/:id`
- **Request Body**
    `-`
- **Return Body**
    ```JSON
    {
    "userID": 6,
    "name": "testuser",
    "email": "blub@blub.com",
    "active": false,
    "isAdmin": false,
    "balance": 0
    }```
    
#### Add Transaction
Adds a transaction.

- **Method**
    `POST`
- **Authentication**
    Bearer-Token with JWT is needed.
    
    Example header:
    ```
    POST /v1/user/6/transaction HTTP/1.1
    Host: localhost:8080
    Content-Type: application/json
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDE4LTAyLTA2VDEzOjEwOjIwLjU0NDI1NTM0OSswMTowMCIsImlzQWRtaW4iOmZhbHNlLCJpc0Jsb2NrZWQiOmZhbHNlLCJ1c2VySUQiOjZ9.TYU_PK9fBf8xZW99CuphrByzDXSVVfV04YCt2oTjWKM
    Cache-Control: no-cache
    ```
- **URL**
    `/v1/user/:id/transaction`
- **Request Body**
    ```JSON
    {
        "userID": 6,
        "sku": 151,
        "tag": "Supercool tag which describes the transaction"
    }
    ```
    or
    ```JSON
    {
        "userID": 6,
        "value": 23.42,
        "tag": "Supercool tag which describes the transaction"
    }
    ```
- **Return Body**
    `-`
    
### Products

| Method      | URL                     | Description                     |
| --------    | --------                | --------                        |
| GET         | /v1/product             | Product index                   |
| POST        | /v1/product             | Add product                     |
| GET         | /v1/product/:id         | Product detail                  |
| PUT         | /v1/product/:id         | TODO: Change product            |
| POST        | /v1/product/:id/stock   | TODO: Change stock +-           |
| DELETE      | /v1/product/:id         | TODO: Delete product            |

#### Product Index
Returns an index of all products with quantity in stock.

- **Method**
    `GET`
- **URL**
    `/v1/product`
- **Request Body**
    `-`
- **Return Body**
    ```json
    [
    {
        "productID": 2,
        "SKU": 151,
        "Name": "test product",
        "GTIN": 123456,
        "price": 12.4,
        "visibility": false,
        "category": null,
        "quantity": 500,
        "quantityUnit": "g",
        "stock": 12
    },
    {
        "productID": 3,
        "SKU": 245,
        "Name": "other product",
        "GTIN": 754544,
        "price": 11.25,
        "visibility": false,
        "category": null,
        "quantity": 0.5,
        "quantityUnit": "liter",
        "stock": 0
    }
    ]
    ```
#### Product Detail
Returns product detail based on SKU.

- **Method**
    `GET`
- **URL**
    `/v1/produkt/:sku`
- **Request Body**
    `-`
- **Return Body**
    ```JSON
    {
    "productID": 2,
    "SKU": 151,
    "Name": "testproduct with new name",
    "GTIN": 123456,
    "price": 12.4,
    "visibility": false,
    "category": null,
    "quantity": 500,
    "quantityUnit": "g",
    "stock": 0
    }```
    

### Categories

| Method      | URL                         | Description                             |
| --------    | --------                    | --------                                |
| GET         | /v1/category                | TODO: Category detail tree              |
| POST        | /v1/category                | TODO: Add category                      |
| GET         | /v1/category/:id            | TODO: Category detail                   |
| PUT         | /v1/category/:id            | TODO: Change cateogry                   |
| POST        | /v1/category/:id/product    | TODO: Add/Delete product from category  |
| DELETE      | /v1/category/:id            | TODO: Delete category                   |


## Licence
This project is licenced under MIT licence. See LICENCE file.
