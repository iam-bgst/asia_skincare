# API For PAKET PRODUCT ASIA SKINCARE

1. [Add](#markdown-header-add-paket)
2. [Get](#markdown-header-get)
3. [Update](#markdown-header-update)
4. [List](#markdown-header-list)
5. [Delete](#markdown-header-delete)
6. [Update Product](#markdown-header0update-product)

## Add

*JSON POST*

``` json
{
    "name":"nama paket",
    "product":["id_product","id_product"],
    "pricing":[
        {
            "membership":"id_membership",
            "price":12345
        }
    ],
    "stock":1234,
    "point":1234,
    "image":"base64(image)"
}
```

*URL*

``` bash
method:POST
ip:port/paket/add
```

*RESPOSE*

``` json
{
    "message": "created",
    "status":  "ok",
}
```

## Get

*URL*

``` bash
method:POST
ip:port/paket/get/:id_paket/:id_membership
```

*RESPONSE JSON*

``` json
{
    "data": {
        "_id": "5097a50a-e429-4429-8c86-fba6432ced58",
        "name": "Paket1",
        "product": [
            {
                "_id": "fb38c192-cc95-466a-ac4e-0f56ab055eeb",
                "name": "Serum Malam",
                "image": ""
            },
            {
                "_id": "572a54d7-2e75-4fe0-9b48-5f899e517acc",
                "name": "Serum Siang",
                "image": ""
            }
        ],
        "pricing": {
            "membership": {
                "_id": "f375391d-c27f-467e-9554-0521a467eb21",
                "name": "Reseller"
            },
            "price": 1000
        },
        "stock": 100,
        "point": 0,
        "image": "/picture/paket/5097a50a-e429-4429-8c86-fba6432ced58.png"
    },
    "status": "ok"
}
```

## Update

*JSON POST*

``` json
{
    "name":"nama paket",
    "product":[],
    "pricing":[],
    "stock":1234,
    "point":1234,
    "image":"base64(image)"
}
```

*URL*

``` bash
method:POST
ip:port/paket/update/:id_paket
```

*JSON RESPONSE*

``` json
{
    	"message": "updated",
		"status":  "ok",
}
```

## List
*URL*

``` bash
method:POST
ip:port/paket/list?membership=id_membership
```
*JSON RESPONSE*
```json
{
    "current_page": 1,
    "data": null,
    "from": 1,
    "last_page": 0,
    "next_page": "",
    "per_page": 5,
    "prev_page": "",
    "status": "Ok",
    "to": 5,
    "total": 0
}
```

## Delete
*URL*

``` bash
method:POST
ip:port/paket/delete/:id_paket
```

*JSON RESPONSE*
```json
{
    "message": "deleted",
	"status":  "ok",
}
```

## Update Product
*JSON POST*
```json
{
    "product":["id_product","id_product"]
}
```

*URL*
``` bash
method:POST
ip:port/paket/update_product/:id_paket
```

*JSON RESPONSE*
```json
    "message": "succes updated",
	"status":  "ok",
```