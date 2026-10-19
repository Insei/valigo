package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo/shared"
)

func codes(errs []shared.Error) []string {
	var codes []string
	for _, err := range errs {
		codes = append(codes, err.Code)
	}
	return codes
}

func TestLimitsFromTheStore(t *testing.T) {
	store := &Store{
		maxTitle: map[string]int{"free": 5, "pro": 50},
		maxTags:  map[string]int{"free": 1, "pro": 5},
		flags:    map[string]bool{},
	}
	v := NewValidator(store)
	post := &Post{Title: "A long title", Tags: []string{"a", "b"}, Priority: 9}
	free, pro := WithPlan(context.Background(), "free"), WithPlan(context.Background(), "pro")

	assert.Equal(t, []string{"string.max_runes", "slice.max_len"}, codes(v.ValidateTyped(free, post)))
	assert.Nil(t, v.ValidateTyped(pro, post))

	store.SetMaxTitle("pro", 3)
	errs := v.ValidateTyped(pro, post)
	assert.Equal(t, []string{"string.max_runes"}, codes(errs))
	assert.Equal(t, map[string]any{"max": 3}, errs[0].Params)
	assert.Equal(t, "Cannot be longer than 3 characters", errs[0].Message)

	store.flags["priorities"] = true
	assert.Equal(t, []string{"string.max_runes", "number.interval"}, codes(v.ValidateTyped(pro, post)))
}
