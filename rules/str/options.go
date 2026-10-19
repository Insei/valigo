package str

// regexpOptions is a struct that represents the options for a regular expression.
type regexpOptions struct {
	localeKey string
}

// RegexpOption is an interface that represents an option for a regular expression.
type RegexpOption interface {
	apply(*regexpOptions)
}

// regexpOptionFunc is a function type that implements the RegexpOption interface.
type regexpOptionFunc func(*regexpOptions)

// apply applies the regexpOptionFunc to the given regexpOptions.
func (f regexpOptionFunc) apply(r *regexpOptions) {
	f(r)
}

// WithRegexpLocaleKey returns a RegexpOption that sets the localeKey field of the regexpOptions.
// It takes a locale key string as an argument.
func WithRegexpLocaleKey(localeKey string) RegexpOption {
	return regexpOptionFunc(func(o *regexpOptions) {
		if len(localeKey) > 0 {
			o.localeKey = localeKey
		}
	})
}
