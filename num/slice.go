package num

import (
	"github.com/insei/valigo/shared"
)

type StringSliceFieldConfigurator struct {
	*shared.SliceFieldConfigurator
}

func (c *StringSliceFieldConfigurator) Trim() *StringSliceFieldConfigurator {
	c.Append(func(v []*any) bool {
		for _, s := range v {
			if s == nil {
				continue
			}

			//*s = strings.TrimSpace(*s)
		}
		return true
	}, "")
	return c
}
