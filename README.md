# Hack-n-Pay
Ein Tool für Hackspaces, Vereine, Gruppen und Firmen um eine Strichliste digital zu verwalten.
## API Dokumentation

### User

| Methode     | URL                         | Bezeichnung                 |
| --------    | --------                    | --------                    |
| GET         | /v1/user                    | User Index                  |
| POST        | /v1/user                    | Sign Up                     |
| POST        | /v1/login                   | Login                       |
| GET         | /v1/user/:id                | User Detail                 |
| POST        | /v1/user/:id/transaction    | Kontostand ändern +-        |
| PUT         | /v1/user/:id                | TODO: Benutzerdaten ändern  |
| DELETE      | /v1/user/:id                | TODO: Benutzer Löschen      |

#### User Index
Gibt einen Index aller User zurück.

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
Fügt Benutzer hinzu

- **Method**
    `POST`
- **URL**
    `/v1/user`
- **Request Body**
    ```json
    {
        "name": "username",
        "email": "blub@blub.de",
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
Login und erstellen eines JWT

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
   JWT-Token. Beispiel: `"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOiIyMDE4LTAxLTI5VDExOjE0OjE5LjIzNjAwMjE2KzAxOjAwIiwiSXNBY3RpdmUiOmZhbHNlLCJJc0FkbWluIjpmYWxzZSwiTmFtZSI6InRlc3R1c2VyIiwidXNlcklEIjo2fQ.Vm2NKt8a5KXKQhEaeb1wBQbPplAXlrkhZ05ZgaKHIAY"`

#### User Detail
Gibt Userdetails anhand der ID zurück

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
    "email": "blub@blub.de",
    "active": false,
    "isAdmin": false,
    "balance": 0
    }```
    
#### Add Transaction
Fügt eine Transaktion hinzu.

- **Method**
    `POST`
- **Authentication**
    Bearer-Token mit JWT wird benötigt.
    
    Beispiel für einen Header:
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
        "SKU": 151,
        "value": 23.42,
        "tag": "Supercool tag which describes the transaction"
    }
    ```
- **Return Body**
    `-`

### Produkte

| Methode     | URL                     | Bezeichnung                     |
| --------    | --------                | --------                        |
| GET         | /v1/product             | Produkt Index                   |
| POST        | /v1/product             | Produkt hinzufügen              |
| GET         | /v1/product/:id         | Produkt Detail                  |
| PUT         | /v1/product/:id         | TODO: Produkt ändern            |
| POST        | /v1/product/:id/stock   | TODO: Lagerbestand ändern +-    |
| DELETE      | /v1/product/:id         | TODO: Produkt Löschen           |

#### Produkt Index
Gibt einen Index aller Produkte zurück.

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
        "Name": "testprodukt neuer name",
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
        "Name": "ganz anderes produkt",
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
#### Produkt Detail
Gibt Produktdetails anhand der SKU zurück

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
    "Name": "testprodukt neuer name",
    "GTIN": 123456,
    "price": 12.4,
    "visibility": false,
    "category": null,
    "quantity": 500,
    "quantityUnit": "g",
    "stock": 0
    }```

## Licence
This project is licenced under MIT licence. See LICENCE file.
