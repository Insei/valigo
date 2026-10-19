package shared

import "fmt"

// Error define a custom error struct.
type Error struct {
	// A human-readable error message.
	Message string
	// The location where the error occurred (e.g., file and line number).
	Location string
	// Additional error information (e.g., a value that caused the error).
	Value any
	// Code is the stable code of the rule, e.g. string.max_len, for the clients that translate the messages themselves.
	Code string `json:",omitempty"`
	// Params are the named parameters of the message, e.g. {"max": 64}, they are shared by the errors of the rule.
	Params map[string]any `json:",omitempty"`
}

// Error implements the error interface by defining an Error() method.
func (e Error) Error() string {
	if e.Location == "" && e.Value == nil {
		return e.Message
	}
	return fmt.Sprintf("%s (%s: %v)", e.Message, e.Location, e.Value)
}
