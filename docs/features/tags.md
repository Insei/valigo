# Rules from tags

`b.Tags("validate")` configures the rules of the tags of the fields of the struct and of its nested structs in the syntax of [go-playground/validator](https://github.com/go-playground/validator): `validate:"required,max=64"`. It is for a migration: the structs keep their tags, and the rules move to code one by one — the tags and the rules in code are mixed in one `Configure`.

The rules of the tags are the rules of valigo: the same messages in the languages of the users, the [codes](errors.md#codes-for-clients), the [JSON Schema](schema.md) and the speed — a valid object is validated without allocations, see [Performance](../performance.md).

## Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/insei/valigo"
	"github.com/insei/valigo/str"
)

type Address struct {
	City string `validate:"required,max=32"`
}

type SignUp struct {
	Email    string   `validate:"required,email"`
	Password string   `validate:"required,min=8"`
	Confirm  string   `validate:"eqfield=Password"`
	Age      int      `validate:"omitempty,gte=18,lte=120"`
	Tags     []string `validate:"max=3,dive,required,max=16"`
	Address  *Address // the rules of City apply when the address is present
	Nick     string
}

func main() {
	v := valigo.New()
	valigo.Configure[SignUp](v, func(b valigo.Configurator[SignUp], s *SignUp) {
		b.Tags("validate")
		// the rules in code with the tags
		b.Field(&s.Nick).As(str.Of).Trim().Alphanumeric()
	})

	errs := v.ValidateTyped(context.Background(), &SignUp{
		Email: "ann", Password: "secret", Confirm: "secrets", Age: 16, Tags: []string{"go", ""}, Nick: " a b ",
	})
	for _, err := range errs {
		fmt.Println(err.Location, "-", err.Message, "-", err.Code)
	}
}
```

Output:

```text
Email - Should be email address - string.email
Password - Cannot be shorter than 8 characters - string.min_runes
Confirm - Should be equal to Password - string.equal_field
Age - Cannot be less than 18 - number.min
Tags[1] - Should be fulfilled - string.required
Nick - Should contain only letters and digits - string.alphanumeric
```

## Rules

| Fields | Rules |
|---|---|
| all | `required`, `omitempty` — the next rules are skipped for the zero value |
| strings | `min`, `max`, `len` in characters, `eq`, `oneof=a b`, `startswith`, `endswith`, `contains`, `email`, `url`, `http_url`, `uuid`, `ip`, `ipv4`, `ipv6`, `hostname`, `hostname_rfc1123`, `alphanum`, `ascii`, `lowercase`, `uppercase`, `e164`, `base64`, `base64url`, `json`, `semver`, `cidr`, `mac`, `hexcolor`, `iso3166_1_alpha2`, `iso4217`, `credit_card`, `luhn_checksum`, `datetime=2006-01-02`, `eqfield` |
| numbers | `min`, `gte`, `max`, `lte`, `eq`, `oneof=1 2`, `latitude`, `longitude`, `gtefield`, `ltefield` |
| `bool` | `required`: `true` |
| `time.Time` | `required`, `gtfield`, `ltfield` |
| slices | `required`, `min`, `max`, `len` of the length, `unique`, `dive` — the next rules are the rules of the elements |
| maps | `required`, `min`, `max`, `len` of the length, `dive` — the next rules are the rules of the values of `map[string]V` |
| pointers to nested structs | `required`; the fields have their own tags |

Pointers to the values have the rules of the values, a nil pointer is absent. `eqfield`, `gtefield` and the other rules on fields take the name of a field of the same struct, or its nested field `Inner.Field`.

## Differences from go-playground/validator

- `required` is `Required` of valigo: a string is not empty, a number is not zero, a slice and a map have elements — go-playground/validator accepts an empty non-nil slice.
- The fields of a nil nested struct are skipped, like by go-playground/validator. The rules in code check `Required` of these fields, see [Absent values](absent.md).
- `dive` into the elements of structs, `keys`/`endkeys`, the alternatives `|`, `gt`/`lt` of numbers and the custom tags of go-playground/validator aren't supported: write these rules in code, with `When` and `Custom`.
- An unknown rule panics in `Configure`, like an invalid rule in code, so a typo doesn't turn into a missing check.
- The rules after `omitempty` aren't in the JSON Schema: they don't apply to every value.
