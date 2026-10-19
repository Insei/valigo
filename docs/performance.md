# Performance

The [`bench`](https://github.com/Insei/valigo/tree/main/bench) module compares valigo with [go-playground/validator](https://github.com/go-playground/validator) and [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) on the same struct and the same rules: a required name up to 64 characters, a required email, an age in `[18, 120]`, up to 5 required tags up to 16 characters and a required city of a nested struct. A test checks that the libraries find the same errors.

```go
b.Field(&u.Name).As(str.Of).Required().MaxRunes(64)
b.Field(&u.Email).As(str.Of).Required().Email()
b.Field(&u.Age).As(num.Of).AnyOfInterval(18, 120)
b.Slice(&u.Tags).As(str.Slice).MaxLen(5).Each(func(tag *str.BaseConfigurator) { tag.Required().MaxRunes(16) })
b.Field(&u.Address.City).As(str.Of).Required()
```

## Results

Intel Core i9-14900HX, Go 1.27.0, one CPU (`-cpu 1`), the median of 6 runs.

A valid object:

| Library | ns/op | B/op | allocs/op |
|---|---|---|---|
| valigo | **485** | **0** | **0** |
| go-playground/validator | 1 830 | 296 | 16 |
| ozzo-validation | 2 900 | 2 000 | 38 |

An invalid object with 6 errors:

| Library | ns/op | B/op | allocs/op |
|---|---|---|---|
| valigo, translated messages | 3 620 | 3 152 | 40 |
| go-playground/validator, without messages | 2 630 | 1 688 | 45 |
| go-playground/validator, with `Error()` of every error | 4 120 | 2 488 | 69 |
| ozzo-validation, with `Error()` | 12 700 | 11 026 | 140 |

- valigo builds the rules once in `Configure` and reads the fields by precomputed accessors, so valid objects cost no allocations; go-playground/validator parses the tags once too, but works through `reflect` on every call, ozzo-validation builds its rules on every call of `Validate`.
- valigo builds the translated messages of the errors at once, go-playground/validator returns the errors without messages and builds them by `Error` or `Translate`: compare valigo with the row "with `Error()`".

## Run

```bash
cd bench
go test -run XXX -bench . -benchmem -cpu 1 -count 6
```

Every pull request gets a comment with the benchmarks of the core compared with the target branch by benchstat, the `bench-diff` workflow.
