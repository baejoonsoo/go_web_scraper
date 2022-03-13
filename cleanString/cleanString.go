package cleanstring

import "strings"

// CleanString cleans a string
func CleanString(str string) string{
	return strings.Join(strings.Fields(strings.TrimSpace(str))," ")
}