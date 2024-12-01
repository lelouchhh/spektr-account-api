package middleware

import "testing"

func TestContainsForbiddenChars(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		// Valid passwords
		{"Valid password with no forbidden chars", "Str0ngP@ssword", false},
		{"Valid password with special chars", "C0mpl3xP@ss!", false},

		// Invalid passwords
		{"Invalid password with forbidden char /", "Invalid/Pass", true},
		{"Invalid password with forbidden char ;", "Invalid;Pass", true},
		{"Invalid password with forbidden char <", "Invalid<Pass", true},
		{"Invalid password with multiple forbidden chars", "Invalid<|>Pass", true},

		// Valid email addresses
		{"Valid email address", "user@example.com", false},
		{"Valid email with numbers", "user123@example.com", false},

		// Invalid email addresses
		{"Email with forbidden char <", "user<test@example.com", true},
		{"Email with forbidden char >", "user>test@example.com", true},
		{"Email with forbidden char |", "user|test@example.com", true},

		// Valid phone numbers
		{"Valid phone number", "1234567890", false},
		{"Phone number with country code", "+1234567890", false},

		// Invalid phone numbers
		{"Phone number with forbidden char <", "<1234567890", true},
		{"Phone number with forbidden char >", "1234567890>", true},
		{"Phone number with forbidden char &", "1234&567890", true},

		// Edge cases
		{"Empty string", "", false},
		{"Single forbidden char", "/", true},
		{"Single non-forbidden char", "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsForbiddenChars(tt.password); got != tt.want {
				t.Errorf("ContainsForbiddenChars(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}
