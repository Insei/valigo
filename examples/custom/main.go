package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/insei/valigo"
	"github.com/insei/valigo/shared"
)

type User struct {
	Name     string
	LastName string
}

func main() {
	v := valigo.New() //v := valigo.New()
	valigo.Configure[User](v, func(builder valigo.Builder[User], obj *User) {
		// Custom on struct type
		builder.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *User) []shared.Error {
			if obj.Name != "Rebecca" {
				format := "Only Rebecca name is allowed" // you can add translations if you want, see translations example
				return []shared.Error{h.ErrorT(ctx, &obj.Name, obj.Name, format)}
			}
			return nil
		})
		//Custom on field
		builder.String(&obj.LastName).
			Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *string) []shared.Error {
				if *value != "Rebecca" {
					localeKey := "Only \"Rebecca\" is allowed" // you can add translations if you want, see translations example
					return []shared.Error{h.ErrorT(ctx, value, localeKey)}
				}
				return nil
			})
	})
	sender := &User{
		Name: "John",
	}
	errs := v.Validate(context.Background(), sender)
	errsJson, _ := json.Marshal(errs)
	fmt.Print(string(errsJson))
}
