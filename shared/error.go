package shared

import "fmt"

type Error struct {
	Message  string
	Location string
	Value    any
}

func (e Error) Error() string {
	if e.Location == "" && e.Value == nil {
		return e.Message
	}
	return fmt.Sprintf("%s (%s: %v)", e.Message, e.Location, e.Value)
}
