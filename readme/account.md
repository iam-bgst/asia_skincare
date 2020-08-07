# API For ACCOUNT ASIA SKINCARE

1. [Register](#markdown-header-register)
2. [Check Account](#markdown-header-check-account)

## Register
*POST JSON*

```json
{
    "name":"nama saya",
    "email":"nama_saya@gmail.com",
    "phonenumber":"0812312",
    "address":"kediri",
    "membership":"657d0c2e-625f-4121-a341-a023d8941568"
}
```

*NB: 657d0c2e-625f-4121-a341-a023d8941568 membership khusus reseller, sementara belum ada list membership*

*URL*
```bash
method:POST
ip:port/account/register
```

*RESPOSE*
```json
{
    "message":"Registed",
    "status": "Ok"
}
```

## Check Account
*URL*
```bash
method:POST
ip:port/account/checkaccount?phone=12345
```

*NB : Query phone harus diisi*

*RESPOSE*
```json
{
    "data": {
        "_id": "belum_diketahui",
        "name": "nama saya",
        "email": "nama_saya@gmail.com",
        "phonenumber": 812312,
        "point": 0,
        "address": "kediri",
        "confirmcode": 0,
        "membership": {
            "_id": "657d0c2e-625f-4121-a341-a023d8941568",
            "name": "Reseller"
        },
        "image": "",
        "status": "active"
    },
    "status": "ok"
}
```