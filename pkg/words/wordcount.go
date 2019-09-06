package words

import "strings"

// Count counts number of words in a string
func Count(s *string) int {
	words := strings.Fields(*s)
	return len(words)
}
