# <img height="32" src="https://raw.githubusercontent.com/itrabbit/sieve/master/logo.png">

GoLang Package for filtering fields models during serialization

**Support:**

- Scopes filtration fields structures
- Exclude (when comparing field values)
- Export values of higher level

**Plans**:
- Add versioning for fields
- Add code generation for fast serialization
- Add support XML, YML serializers

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
    UpdatedAt uint64   `json:"updated_at" sieve:"eef:CreatedAt"`
    Options   []Option `json:"options" sieve:"ek:Name"`
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