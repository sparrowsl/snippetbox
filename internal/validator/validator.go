package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// Returns true if the fieldError map is empty - no entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Adds an error to the fieldError map,
// replaces the given field value if it exists
func (v *Validator) AddFieldError(key string, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, ok := v.FieldErrors[key]; !ok {
		v.FieldErrors[key] = message
	}
}

// Adds an error message to the fieldError map only if a check is not 'ok'
func (v *Validator) CheckField(ok bool, key string, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// Returns true if a value is in a list of permitted integers
func PermittedInt(n int, values ...int) bool {
	for i := range values {
		if n == values[i] {
			return true
		}
	}

	return false
}
