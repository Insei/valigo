// Package problem writes validation errors as the problem details of RFC 9457 (RFC 7807 before it):
// application/problem+json with the errors of the fields.
//
//	v := valigo.New(valigo.WithErrorsTransformer(problem.Transformer()))
//	if errs := v.Validate(ctx, req); errs != nil {
//		problem.From(errs).Write(w)
//	}
package problem

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/insei/valigo/shared"
)

// ContentType is the media type of the problem details.
const ContentType = "application/problem+json"

// Details is the problem details of RFC 9457 with the errors of the fields in the errors member.
type Details struct {
	// Type is the URI of the type of the problem, about:blank when it is empty.
	Type string `json:"type,omitempty"`
	// Title is a short summary of the problem.
	Title string `json:"title"`
	// Status is the HTTP status code.
	Status int `json:"status"`
	// Detail is an explanation of this occurrence of the problem.
	Detail string `json:"detail,omitempty"`
	// Instance is the URI of this occurrence of the problem.
	Instance string `json:"instance,omitempty"`
	// Errors are the errors of the fields.
	Errors []FieldError `json:"errors"`
}

// FieldError is an error of a field.
type FieldError struct {
	// Detail is the translated message.
	Detail string `json:"detail"`
	// Pointer is the JSON Pointer of the field in the request, e.g. #/tags/1, built from the location.
	Pointer string `json:"pointer"`
	// Location is the location of the field, e.g. tags[1].
	Location string `json:"location"`
	// Code is the stable code of the rule, e.g. string.max_len.
	Code string `json:"code,omitempty"`
	// Params are the named parameters of the message, e.g. {"max": 64}.
	Params map[string]any `json:"params,omitempty"`
}

// Option configures the problem details.
type Option func(d *Details)

// WithType sets the URI of the type of the problem.
func WithType(uri string) Option {
	return func(d *Details) { d.Type = uri }
}

// WithTitle sets the title, "Validation failed" by default.
func WithTitle(title string) Option {
	return func(d *Details) { d.Title = title }
}

// WithStatus sets the HTTP status code, 422 Unprocessable Content by default.
func WithStatus(status int) Option {
	return func(d *Details) { d.Status = status }
}

// WithDetail sets the explanation of the problem.
func WithDetail(detail string) Option {
	return func(d *Details) { d.Detail = detail }
}

// New creates the problem details of the validation errors.
func New(errs []shared.Error, opts ...Option) *Details {
	d := &Details{Title: "Validation failed", Status: http.StatusUnprocessableEntity, Errors: make([]FieldError, len(errs))}
	for i, err := range errs {
		d.Errors[i] = FieldError{Detail: err.Message, Pointer: Pointer(err.Location), Location: err.Location, Code: err.Code, Params: err.Params}
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Transformer returns the errors transformer of valigo.WithErrorsTransformer:
// Validate returns one error, the *Details of all errors.
func Transformer(opts ...Option) func(errs []shared.Error) []error {
	return func(errs []shared.Error) []error {
		return []error{New(errs, opts...)}
	}
}

// From returns the problem details of the errors of Validate with the Transformer, nil if there are none.
func From(errs []error) *Details {
	for _, err := range errs {
		var d *Details
		if errors.As(err, &d) {
			return d
		}
	}
	return nil
}

// Error returns the title and the first error of a field.
func (d *Details) Error() string {
	switch len(d.Errors) {
	case 0:
		return d.Title
	case 1:
		return d.Title + ": " + d.Errors[0].Location + ": " + d.Errors[0].Detail
	}
	return d.Title + ": " + d.Errors[0].Location + ": " + d.Errors[0].Detail + " and other errors"
}

// pointerEscaper escapes ~ and / of the reference tokens of a JSON Pointer by RFC 6901.
var pointerEscaper = strings.NewReplacer("~", "~0", "/", "~1")

// Write writes the problem details to the response with the status and the application/problem+json content type.
func (d *Details) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(d.Status)
	return json.NewEncoder(w).Encode(d)
}

// Pointer converts a location of an error, e.g. address.city or tags[1], to a JSON Pointer
// of RFC 6901 in a URI fragment: #/address/city, #/tags/1.
func Pointer(location string) string {
	var b strings.Builder
	b.WriteString("#")
	for part := range strings.FieldsFuncSeq(location, func(r rune) bool { return r == '.' || r == '[' || r == ']' }) {
		b.WriteByte('/')
		b.WriteString(pointerEscaper.Replace(part))
	}
	return b.String()
}
