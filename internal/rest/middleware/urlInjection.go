package middleware

import (
	"regexp"
)

func IsSuspicious(input string) bool {
	// Regular expression to catch dangerous characters or patterns
	re := regexp.MustCompile(`(?:[;'"<>()%]|--|\/\*|\*\/|\b(select|union|drop|update|insert)\b|[^a-zA-Z0-9@!#$%&*])`)
	return re.MatchString(input)
}
