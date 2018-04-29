# <img height="48" src="https://raw.githubusercontent.com/itrabbit/sieve/master/logo.png">

[![Build Status](https://travis-ci.org/itrabbit/sieve.svg?branch=master)](https://travis-ci.org/itrabbit/sieve)
 [![CodeCov](https://codecov.io/gh/itrabbit/sieve/branch/master/graph/badge.svg)](https://codecov.io/gh/itrabbit/sieve)
 [![GoDoc](https://godoc.org/github.com/itrabbit/sieve?status.svg)](https://godoc.org/github.com/itrabbit/sieve)

GoLang Package for filtering fields models during serialization

**Support:**

- Scopes filtration fields structures
- Exclude strategy (equal fieldValue, equal value)
- Export values of higher level

**Plans**:
- Add versioning for fields
- Add code generation for fast serialization
- Add support XML, YML serializers


### Tag options

- `s`, `scopes` - for the pane fields in structures;

- `k`, `exportKeys` - exported field or value from the structure fields;

- `ef`, `eef`, `excludeEqualField` - exclusion rule field (comparison with another field);

- `ev`, `eev`, `excludeEqualValue` -  exclusion rule field (comparison with value);

- `e.any` - if you are using multiple conditions, exceptions, if specified, is triggered when any condition.

**Example:**
```go
type A struct {
    Status string `sieve:"ev:unknown;ev:null;e.any"

    CreatedAt time.Time `sieve:"s:private"`
    UpdatedAt time.Time `sieve:"s:private,admin;ef:CreatedAt"`
}
```

## Example

```go

...

type Option struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}

type Object struct {
    idx       uint64
    Name      string   `json:"name" sieve:"s:*"`
    FullName  string   `json:"full_name" sieve:"s:private"`
    CreatedAt uint64   `json:"created_at"`
    UpdatedAt uint64   `json:"updated_at" sieve:"ef:CreatedAt"`
    Options   []Option `json:"options" sieve:"k:Name"`
}

...

func main() {
    obj := Object{
        100,
        "One",
        "Full",
        100,
        100,
        []Option{
            Option{"1", "Cool"},
            Option{"2", "Great"},
            Option{"3", "Best"},
        },
    }
    if b, err := json.Marshal(Sieve(&obj, "public")); err == nil {
        fmt.Println(string(b))
        // -> {"created_at":100,"name":"One","options":["1","2","3"]}
    }
    if b, err := json.Marshal(Sieve(&obj, "private")); err == nil {
        fmt.Println(string(b))
        // -> {"created_at":100,"full_name":"Full","name":"One","options":["1","2","3"]}
    }
}
```


# Donation to development

`BTC: 1497z5VaY3AUEUYURS5b5fUTehVwv7wosX`

`DASH: XjBr7sqaCch4Lo1A7BctQz3HzRjybfpx2c`

`XRP: rEQwgdCr8Jma3GY2s55SLoZq2jmqmWUBDY`

`PayPal / Yandex Money: garin1221@yandex.ru`


## License

JUST is licensed under the [MIT](LICENSE).