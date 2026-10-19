# Performance

The [`bench`](https://github.com/Insei/valigo/tree/main/bench) module compares valigo with [go-playground/validator](https://github.com/go-playground/validator), [ozzo-validation](https://github.com/go-ozzo/ozzo-validation), [zog](https://github.com/Oudwins/zog), [gookit/validate](https://github.com/gookit/validate) and hand-written Go code on the same struct and the same rules: a required name up to 64 characters, a required email, an age in `[18, 120]`, up to 5 required tags up to 16 characters and a required city of a nested struct. A test checks that the libraries find the same errors.

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

| Library                            | ns/op   | B/op  | allocs/op |
| ---------------------------------- | ------- | ----- | --------- |
| hand-written code, the lower bound | 286     | 0     | 0         |
| valigo                             | **407** | **0** | **0**     |
| zog                                | 1 690   | 6     | 2         |
| go-playground/validator            | 1 785   | 296   | 16        |
| ozzo-validation                    | 2 950   | 2 000 | 38        |
| gookit/validate                    | 6 880   | 4 560 | 65        |

An invalid object with 6 errors:

| Library                                                | ns/op     | B/op      | allocs/op |
| ------------------------------------------------------ | --------- | --------- | --------- |
| hand-written code, errors without messages             | 590       | 288       | 8         |
| valigo, translated messages                            | **2 205** | **1 168** | **14**    |
| go-playground/validator, without messages              | 2 665     | 1 688     | 45        |
| go-playground/validator, with `Error()` of every error | 4 150     | 2 488     | 69        |
| zog, issues with messages                              | 4 600     | 1 056     | 24        |
| ozzo-validation, with `Error()`                        | 12 665    | 11 026    | 140       |
| gookit/validate                                        | 13 250    | 9 569     | 120       |

valigo is about 1.4 times slower than the same rules written by hand on a valid object, it pays for the closures of the rules; it doesn't allocate either. The fields are read by the offsets of fmap without reflect, following the pointers of nested structs. Code generators of validators produce code like the hand-written one, so they are expected to be about as fast as it.

- valigo builds the rules once in `Configure` and reads the fields by precomputed accessors, so valid objects cost no allocations; go-playground/validator parses the tags once too, but works through `reflect` on every call, ozzo-validation builds its rules on every call of `Validate`.
- gookit/validate creates a validation per call, it is its API, and reports the errors of the elements of a slice under one key `Tags.*`; zog builds the messages of its issues.
- valigo builds the translated messages of the errors at once, go-playground/validator returns the errors without messages and builds them by `Error` or `Translate`: compare valigo with the row "with `Error()`".
- The rules of valigo append the errors to one pooled buffer, `Validate` copies it once; most allocations of an invalid object are the messages and the boxed values of the errors.

## Run

```bash
cd bench
go test -run XXX -bench . -benchmem -cpu 1 -count 6
```

Every pull request gets a comment with the benchmarks of the core compared with the target branch by benchstat, the `bench-diff` workflow.
