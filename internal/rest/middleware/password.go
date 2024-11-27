package middleware

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidatePassword checks if the password meets required criteria
func ValidatePassword(password string) error {
	// Check for minimum length of 8 characters
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	//
	//// Check for at least one uppercase letter
	//if !hasUpperCase(password) {
	//	return fmt.Errorf("password must contain at least one uppercase letter")
	//}
	//
	//// Check for at least one lowercase letter
	//if !hasLowerCase(password) {
	//	return fmt.Errorf("password must contain at least one lowercase letter")
	//}
	//
	//// Check for at least one digit
	//if !hasDigit(password) {
	//	return fmt.Errorf("password must contain at least one digit")
	//}
	//
	//// Check for at least one special character
	//if !hasSpecialChar(password) {
	//	return fmt.Errorf("password must contain at least one special character")
	//}

	// Check for unwanted characters (e.g., /, ;, etc.)
	if containsForbiddenChars(password) {
		return fmt.Errorf("password contains forbidden characters like /, ;, etc.")
	}
	//
	//// Check if the password contains a common or weak word
	//if containsCommonPassword(password) {
	//	return fmt.Errorf("password is too common or weak")
	//}

	return nil
}

// Helper function to check if the string contains an uppercase letter
func hasUpperCase(s string) bool {
	re := regexp.MustCompile(`[A-Z]`)
	return re.MatchString(s)
}

// Helper function to check if the string contains a lowercase letter
func hasLowerCase(s string) bool {
	re := regexp.MustCompile(`[a-z]`)
	return re.MatchString(s)
}

// Helper function to check if the string contains a digit
func hasDigit(s string) bool {
	re := regexp.MustCompile(`[0-9]`)
	return re.MatchString(s)
}

// Helper function to check if the string contains a special character
func hasSpecialChar(s string) bool {
	re := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	return re.MatchString(s)
}

// Helper function to check if the password contains forbidden characters like /, ;
func containsForbiddenChars(password string) bool {
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

// Helper function to check if the password contains any common or weak words
func containsCommonPassword(password string) bool {
	// Convert password to lowercase for comparison
	lowerPassword := strings.ToLower(password)

	// List of common passwords or weak password patterns
	weakPasswords := []string{
		"password", "123456", "12345", "12345678", "qwerty", "abc123", "letmein", "welcome", "123qwe", "password1",
	}

	// Check if password contains any common weak password
	for _, weak := range weakPasswords {
		if strings.Contains(lowerPassword, weak) {
			return true
		}
	}
	return false
}
