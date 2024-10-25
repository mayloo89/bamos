package forms

type errors map[string][]string

// Add adds an error message for a given form field
func (e errors) Add(field, msg string) {
	e[field] = append(e[field], msg)
}

// Get returns the first error message from a field
func (e errors) Get(field string) string {
	if errorMsg, exists := e[field]; exists {
		return errorMsg[0]
	}
	return ""
}
