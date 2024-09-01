[![codecov](https://codecov.io/gh/Insei/valigo/graph/badge.svg?token=TK4XN3I2Q6)](https://codecov.io/gh/Insei/valigo)
[![build](https://github.com/Insei/valigo/actions/workflows/go.yml/badge.svg)](https://github.com/Insei/valigo/actions/workflows/go.yml)
[![Goreport](https://goreportcard.com/badge/github.com/insei/valigo)](https://goreportcard.com/report/github.com/insei/valigo)
[![GoDoc](https://godoc.org/github.com/insei/valigo?status.svg)](https://godoc.org/github.com/insei/valigo)
# Valigo 
Valigo is a powerfull, zero allocations validation engine with localizations, conditions and custom validation functions support.
## Features
* Conditinal validation
* Custom validation functions
* Localizations
* Zero allocations
* Configured using pointers to structure fields
## Roadmap
* [x] Zero allocations on valid structs
* [x] On Field Conditional validation
* [x] On Field Custom validation
* [x] On Struct Conditional validation
* [x] On Struct Custom validation
* [x] Error translations
* [x] Strings (MaxLen, MinLen, Required, Regexp Pattern, AnyOf, Custom) and Strings Slices validation
* [x] UUID and UUID Slices validation
* [x] Num Validation (int(8,16,32,64), uint(8,16,32,64))
* [ ] Float validation
* [ ] Other default types validations
* [ ] Create validation rules based on default validations tags
