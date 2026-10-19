// The rules of this example read their limits from an external store on every validation:
// the limits of the plan of the tenant from the context, and a feature flag. Changing the store
// changes the validation without a new Configure.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"unicode/utf8"

	"github.com/insei/valigo"
	"github.com/insei/valigo/rules/num"
	"github.com/insei/valigo/rules/str"
	"github.com/insei/valigo/shared"
)

// Store is the source of the limits, a database or a cache in a real service.
type Store struct {
	mu       sync.RWMutex
	maxTitle map[string]int  // the maximum length of a title by plan
	maxTags  map[string]int  // the maximum number of tags by plan
	flags    map[string]bool // feature flags
}

func (s *Store) MaxTitle(plan string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.maxTitle[plan]
}

func (s *Store) MaxTags(plan string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.maxTags[plan]
}

func (s *Store) Enabled(flag string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.flags[flag]
}

func (s *Store) SetMaxTitle(plan string, n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maxTitle[plan] = n
}

type planKey struct{}

// WithPlan puts the plan of the tenant into the context, a middleware does it in a real service.
func WithPlan(ctx context.Context, plan string) context.Context {
	return context.WithValue(ctx, planKey{}, plan)
}

func planOf(ctx context.Context) string {
	plan, _ := ctx.Value(planKey{}).(string)
	return plan
}

type Post struct {
	Title    string
	Tags     []string
	Priority int
}

// NewValidator configures the rules once, the closures read the store on every validation.
func NewValidator(store *Store) *valigo.Validator {
	v := valigo.New()
	valigo.Configure[Post](v, func(b valigo.Configurator[Post], p *Post) {
		b.Field(&p.Title).As(str.Of).Trim().Required().Custom(func(ctx context.Context, h *shared.FieldCustomHelper, title string) []shared.Error {
			if maxLen := store.MaxTitle(planOf(ctx)); utf8.RuneCountInString(title) > maxLen {
				// the built-in key and code: the message is translated with plural forms, clients get string.max_runes
				msg := shared.M("string.max_runes", "validation:string:Cannot be longer than %d characters", shared.P("max", maxLen))
				return []shared.Error{h.ErrorMsg(ctx, title, msg)}
			}
			return nil
		})
		b.Slice(&p.Tags).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, tags shared.Elems[string]) []shared.Error {
			if maxTags := store.MaxTags(planOf(ctx)); tags.Len() > maxTags {
				msg := shared.M("slice.max_len", "validation:slice:Cannot contain more than %d elements", shared.P("max", maxTags))
				return []shared.Error{h.ErrorMsg(ctx, tags.Len(), msg)}
			}
			return nil
		})
		// the rules of the priority are enabled by a feature flag
		priority := b.When(func(context.Context, *Post) bool { return store.Enabled("priorities") })
		priority.Field(&p.Priority).As(num.Of).AnyOfInterval(1, 5)
	})
	return v
}

func main() {
	store := &Store{
		maxTitle: map[string]int{"free": 10, "pro": 100},
		maxTags:  map[string]int{"free": 1, "pro": 10},
		flags:    map[string]bool{"priorities": true},
	}
	v := NewValidator(store)
	post := &Post{Title: "Limits from a store", Tags: []string{"go", "validation"}, Priority: 9}

	show := func(title string, ctx context.Context) {
		body, _ := json.Marshal(v.Validate(ctx, post))
		fmt.Println(title+":", string(body))
	}
	show("free", WithPlan(context.Background(), "free"))
	show("pro", WithPlan(context.Background(), "pro"))

	// the limits change at runtime, the rules aren't configured again
	store.SetMaxTitle("pro", 5)
	show("pro after the change", WithPlan(context.Background(), "pro"))
}
