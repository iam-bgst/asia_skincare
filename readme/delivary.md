# API For DELIVERY ASIA SKINCARE

1. [List City](#markdown-header-list-city)
2. [Check Ongkir](#markdown-header-check-ongkir)

## List City
*URL*
```bash
method:GET
ip:port/delivery/list
```

*JSON RESPONSE*
```json
{
    "current_page": 1,
    "data": [
        {
            "city_id": "1",
            "province_id": "21",
            "province": "Nanggroe Aceh Darussalam (NAD)",
            "type": "Kabupaten",
            "city_name": "Aceh Barat",
            "postal_code": "23681"
        },
        {
            "city_id": "2",
            "province_id": "21",
            "province": "Nanggroe Aceh Darussalam (NAD)",
            "type": "Kabupaten",
            "city_name": "Aceh Barat Daya",
            "postal_code": "23764"
        },
        {
            "city_id": "3",
            "province_id": "21",
            "province": "Nanggroe Aceh Darussalam (NAD)",
            "type": "Kabupaten",
            "city_name": "Aceh Besar",
            "postal_code": "23951"
        },
        {
            "city_id": "4",
            "province_id": "21",
            "province": "Nanggroe Aceh Darussalam (NAD)",
            "type": "Kabupaten",
            "city_name": "Aceh Jaya",
            "postal_code": "23654"
        },
        {
            "city_id": "5",
            "province_id": "21",
            "province": "Nanggroe Aceh Darussalam (NAD)",
            "type": "Kabupaten",
            "city_name": "Aceh Selatan",
            "postal_code": "23719"
        }
    ],
    "from": 1,
    "last_page": 1,
    "next_page": "",
    "per_page": 5,
    "prev_page": "",
    "status": "Ok",
    "to": 5,
    "total": 501
}
```

## Check Ongkir
*URL*
```bash
method:GET
ip:port/delivery/checkongkir?origin=id_city&destination=id_city&weight=1000
```

*JSON RESPONSE*
```json
{
    "data": [
        {
            "name": "Jalur Nugraha Ekakurir (JNE)",
            "code": "jne",
            "cost": [
                {
                    "name": "OKE",
                    "cost": 7000,
                    "estimate": "2-3"
                },
                {
                    "name": "REG",
                    "cost": 8000,
                    "estimate": "1-2"
                },
                {
                    "name": "YES",
                    "cost": 14000,
                    "estimate": "1-1"
                }
            ]
        },
        {
            "name": "POS Indonesia (POS)",
            "code": "pos",
            "cost": [
                {
                    "name": "OKE",
                    "cost": 7000,
                    "estimate": "2-3"
                },
                {
                    "name": "REG",
                    "cost": 8000,
                    "estimate": "1-2"
                },
                {
                    "name": "YES",
                    "cost": 14000,
                    "estimate": "1-1"
                },
                {
                    "name": "Paket Kilat Khusus",
                    "cost": 8000,
                    "estimate": "2-3 HARI"
                },
                {
                    "name": "Express Next Day Barang",
                    "cost": 13500,
                    "estimate": "1 HARI"
                }
            ]
        },
        {
            "name": "Citra Van Titipan Kilat (TIKI)",
            "code": "tiki",
            "cost": [
                {
                    "name": "OKE",
                    "cost": 7000,
                    "estimate": "2-3"
                },
                {
                    "name": "REG",
                    "cost": 8000,
                    "estimate": "1-2"
                },
                {
                    "name": "YES",
                    "cost": 14000,
                    "estimate": "1-1"
                },
                {
                    "name": "Paket Kilat Khusus",
                    "cost": 8000,
                    "estimate": "2-3 HARI"
                },
                {
                    "name": "Express Next Day Barang",
                    "cost": 13500,
                    "estimate": "1 HARI"
                },
                {
                    "name": "ECO",
                    "cost": 4000,
                    "estimate": "5"
                },
                {
                    "name": "REG",
                    "cost": 7000,
                    "estimate": "3"
                }
            ]
        }
    ],
    "status": "ok"
}
```

