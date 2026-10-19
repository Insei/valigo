package shared

import (
	"context"
	"maps"

	"github.com/insei/fmap/v3"
)

// Param is a named parameter of a message: P("max", 64).
type Param struct {
	Name  string
	Value any
}

// P creates a named parameter of a message.
func P(name string, value any) Param {
	return Param{Name: name, Value: value}
}

// Msg is the message of a rule.
type Msg struct {
	// Code is the stable code of the rule for the clients that translate the messages themselves, e.g. string.max_len.
	Code string
	// Key is the locale key of the translated message, Args are its format arguments.
	Key  string
	Args []any
	// Params are the named parameters of the message, e.g. {"max": 64}, shared by the errors of the rule.
	Params map[string]any
}

// M creates the message of a rule with the stable code, the locale key and the named parameters,
// the values of the parameters are the format arguments of the message in their order.
// Create the messages on the configuration, the errors share them.
func M(code, key string, params ...Param) Msg {
	m := Msg{Code: code, Key: key}
	if len(params) > 0 {
		m.Args = make([]any, len(params))
		m.Params = make(map[string]any, len(params))
		for i, p := range params {
			m.Args[i] = p.Value
			m.Params[p.Name] = p.Value
		}
	}
	return m
}

// with returns a copy of the message with one more parameter, it is the last format argument.
func (m Msg) with(name string, value any) Msg {
	params := make(map[string]any, len(m.Params)+1)
	maps.Copy(params, m.Params)
	params[name] = value
	return Msg{Code: m.Code, Key: m.Key, Args: append(m.Args[:len(m.Args):len(m.Args)], value), Params: params}
}

// Error builds the error of the field with the translated message, the code and the parameters of m.
// The helper of the validator passes the named parameters to the translator, for the sources with named placeholders.
func (m Msg) Error(ctx context.Context, h Helper, field fmap.Field, value any) Error {
	var err Error
	if ph, ok := h.(paramsHelper); ok {
		err = ph.ErrorParams(ctx, field, value, m.Key, m.Args, m.Params)
	} else {
		err = h.ErrorT(ctx, field, value, m.Key, m.Args...)
	}
	err.Code, err.Params = m.Code, m.Params
	return err
}

// paramsHelper is a Helper that passes the named parameters of the messages to the translator.
type paramsHelper interface {
	ErrorParams(ctx context.Context, field fmap.Field, value any, key string, args []any, params map[string]any) Error
}

// key returns the message of a rule configured by a locale key, its code is the key.
func key(localeKey string, args []any) Msg {
	return Msg{Code: localeKey, Key: localeKey, Args: args}
}
