package middleware

import (
	"regexp"
	"strings"
)

func IsSuspicious(input string) bool {
	// Regular expression to catch dangerous characters or patterns
	re := regexp.MustCompile(`(?:[;'"<>()%]|--|\/\*|\*\/|\b(select|union|drop|update|insert)\b|[^a-zA-Z0-9@!#$%&*])`)
	return re.MatchString(input)
}

// Helper function to check if the password contains forbidden characters like /, ;
func ContainsForbiddenChars(password string) bool {
	// List of forbidden characters
	forbiddenChars := []string{"/", ";", "\"", "'", "<", ">", "|", ":", "&"}

	// Check if password contains any forbidden character
	for _, char := range forbiddenChars {
		if strings.Contains(password, char) {
			return true
		}
	}
	return false
}
