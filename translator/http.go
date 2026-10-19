package translator

import (
	"context"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// languagesContextKey represents the key for storing the accept-language
// value in the context.
type languagesContextKey struct {
	Key string
}

// languagesContextKeyVal is the instance of languagesContextKey used
// to store the accept-language value in the context.
var languagesContextKeyVal = &languagesContextKey{
	Key: "accept-language",
}

// Language is a structure representing language and quotient in AcceptLanguage header
type language struct {
	name     string
	quotient float32
}

// SortQuotient is a structure representing list language in AcceptLanguage header
type sortQuotient []language

// Len Swap and Less implement the Interface of the standard sort package
// Used to sort languages by priority based on quotient in AcceptLanguage header
func (q sortQuotient) Len() int           { return len(q) }
func (q sortQuotient) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q sortQuotient) Less(i, j int) bool { return q[i].quotient > q[j].quotient }

// getPriorityLanguages takes an acceptLanguage string and returns a slice of priority languages.
// The acceptLanguage string is a comma-separated list of languages with optional quotient values.
// The priority languages are the languages with the highest quotient values.
func getPriorityLanguages(acceptLanguage string) []string {
	languages := make([]language, 0)
	langs := strings.Split(acceptLanguage, ",")

	for _, notParsedLang := range langs {
		if strings.TrimSpace(notParsedLang) == "" {
			continue
		}
		langWithQ := strings.Split(notParsedLang, ";")
		lang, ok := parseLang(strings.TrimSpace(langWithQ[0]))
		if !ok {
			continue
		}
		q, ok := parseQuotient(langWithQ)
		if !ok {
			continue
		}

		languages = append(languages, language{name: lang, quotient: q})
	}

	sort.Sort(sortQuotient(languages))
	priorityLanguages := make([]string, 0)
	for _, lang := range languages {
		if slices.Contains(priorityLanguages, lang.name) {
			continue
		}
		priorityLanguages = append(priorityLanguages, lang.name)
	}

	return priorityLanguages
}

// parseLang takes a language string and returns the parsed language, i.e. its first two characters in lowercase.
// It reports false for the wildcard and malformed languages.
func parseLang(lang string) (string, bool) {
	if len(lang) < 2 {
		return "", false
	}
	return strings.ToLower(lang[0:2]), true
}

// parseQuotient takes a slice of strings and returns the parsed quotient.
// If the slice has a length of 1, the quotient is set to 1.
// Otherwise, the quotient is parsed from the q parameter, false is reported if it is malformed.
func parseQuotient(parts []string) (float32, bool) {
	for _, param := range parts[1:] {
		key, value, _ := strings.Cut(strings.TrimSpace(param), "=")
		if key != "q" {
			continue
		}
		q, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return 0, false
		}
		return float32(q), true
	}
	return 1, true
}

// NewAcceptLanguageMiddleware accepts the request, retrieves the list of languages from Accept-Language header
// creates new context with languagesContextKeyVal key and sends the request further
func NewAcceptLanguageMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var priorityLanguages []string
			acceptLanguage := r.Header.Get("Accept-Language")
			if acceptLanguage != "" {
				priorityLanguages = getPriorityLanguages(acceptLanguage)
			}
			r = r.WithContext(context.WithValue(r.Context(), languagesContextKeyVal, priorityLanguages))
			next.ServeHTTP(w, r)
		})
	}
}

// GetPreferredLanguagesFromContext gets preferred languages from context.Context.
// SHOULD be used with NewAcceptLanguageMiddleware.
func GetPreferredLanguagesFromContext(ctx context.Context) []string {
	preferredAny := ctx.Value(languagesContextKeyVal)
	if preferredAny == nil {
		return nil
	}
	preferredLanguages, _ := preferredAny.([]string)
	return preferredLanguages
}
