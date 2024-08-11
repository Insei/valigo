package str

type regexpOptions struct {
	localeKey string
}

type RegexpOption interface {
	apply(*regexpOptions)
}

type regexpOptionFunc func(*regexpOptions)

func (f regexpOptionFunc) apply(r *regexpOptions) {
	f(r)
}
func WithRegexpLocaleKey(localeKey string) RegexpOption {
	return regexpOptionFunc(func(o *regexpOptions) {
		if len(localeKey) > 0 {
			o.localeKey = localeKey
		}
	})
}
