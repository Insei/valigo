package valigo

type Options struct {
	storage storage
}

type Option interface {
	apply(*Validator)
}
