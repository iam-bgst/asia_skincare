# API For PRODUCT ASIA SKINCARE

1. [Add Product](#markdown-header-add-product)
2. [List Product By Membership](#markdown-header-list-product-by-membership)
3. [Update Product](#markdown-header-update-product)
4. [Get Product](#markdown-header-get-product)
5. [Update Price per Membership](#markdown-header-update-price-per-membership)
6. [Delete Product](#markdown-header-delete-header)

## Add Product

*POST JSON*

``` json
{
    "name": "Serum Malam",
    "stock": 100000,
    "point": 1000,
    "pricing": [
        {
            "membership": "285340a2-431e-4444-8653-9775311d0f2c",
            "price": 10000
        }, {
            "membership": "657d0c2e-625f-4121-a341-a023d8941568",
            "price": 10000
        },{
            "membership": "6ea89ef7-9a3e-4c4f-ac98-196d787a887f",
            "price": 10000
        },{
            "membership": "48383414-ba75-4b54-94e1-6d7909edc8bc",
            "price": 10000
        }
    ]
}
```

*URL*

``` bash
method:POST
ip:port/product/add
```

*RESPOSE*

``` json
{
    "message": "Created",
    "status": "ok"
}
```

## List Product By Membership

*URL*

``` bash
method:GET
ip:port/product?membership=285340a2-431e-4444-8653-9775311d0f2c
```

*NB : Query Membership is required*

*JSON RESPONSE*

``` json
{
    "current_page": 1,
    "data": [
        {
            "_id": "2c2505cd-c2fa-4f28-b971-603716ecdac9",
            "name": "Serum",
            "pricing": {
                "membership": {
                    "_id": "657d0c2e-625f-4121-a341-a023d8941568",
                    "name": "Reseller"
                },
                "price": 10000
            },
            "strock": 1000,
            "point": 10,
            "image": ""
        },
        {
            "_id": "6a814f95-4e0d-465e-b1bc-b3902a91d768",
            "name": "Serum Malam",
            "pricing": {
                "membership": {
                    "_id": "657d0c2e-625f-4121-a341-a023d8941568",
                    "name": "Reseller"
                },
                "price": 10000
            },
            "strock": 100000,
            "point": 1000,
            "image": ""
        }
    ],
    "from": 1,
    "last_page": 1,
    "next_page": "",
    "per_page": 5,
    "prev_page": "",
    "status": "Ok",
    "to": 5,
    "total": 2
}
```

## Update Product

*JSON POST*

``` json
{
    "name": "Serum Malam",
    "stock": 100000,
    "point": 1000,
}
```

*URL*
```bash
method:PUT
ip:port/product/update/id_product
```

*JSON RESPONSE*

```json
{
    "message": "success update product",
    "status": "ok"
}
```
