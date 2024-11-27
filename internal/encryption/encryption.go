package encryption

import (
	"errors"
	"strings"
	"unicode"
)

// Encrypt encrypts an alphanumeric string using a Vigenère cipher with the given key.
func Encrypt(input, key string) (string, error) {
	if !isAlphaNumeric(input) || !isAlphaNumeric(key) {
		return "", errors.New("input and key must contain only numbers and alphabets")
	}

	var encrypted strings.Builder
	keyIndex := 0
	keyLength := len(key)

	for _, char := range input {
		shift := calculateShift(rune(key[keyIndex%keyLength]))
		encrypted.WriteRune(shiftChar(char, shift))
		keyIndex++
	}

	return encrypted.String(), nil
}

// Decrypt decrypts an alphanumeric string encrypted using a Vigenère cipher with the given key.
func Decrypt(input, key string) (string, error) {
	if !isAlphaNumeric(input) || !isAlphaNumeric(key) {
		return "", errors.New("input and key must contain only numbers and alphabets")
	}

	var decrypted strings.Builder
	keyIndex := 0
	keyLength := len(key)

	for _, char := range input {
		shift := calculateShift(rune(key[keyIndex%keyLength]))
		decrypted.WriteRune(shiftChar(char, -shift))
		keyIndex++
	}

	return decrypted.String(), nil
}

// calculateShift calculates the shift value based on a character in the key.
func calculateShift(keyChar rune) int {
	if unicode.IsLower(keyChar) {
		return int(keyChar - 'a')
	} else if unicode.IsUpper(keyChar) {
		return int(keyChar - 'A')
	} else if unicode.IsDigit(keyChar) {
		return int(keyChar - '0')
	}
	return 0
}

// shiftChar shifts a character by the given shift value, wrapping around alphabets and numbers.
func shiftChar(char rune, shift int) rune {
	if char >= 'a' && char <= 'z' {
		return 'a' + (char-'a'+rune(shift)+26)%26
	} else if char >= 'A' && char <= 'Z' {
		return 'A' + (char-'A'+rune(shift)+26)%26
	} else if char >= '0' && char <= '9' {
		return '0' + (char-'0'+rune(shift)+10)%10
	}
	return char
}

// isAlphaNumeric checks if a string contains only numbers and alphabets.
func isAlphaNumeric(input string) bool {
	for _, char := range input {
		if !(char >= 'a' && char <= 'z') &&
			!(char >= 'A' && char <= 'Z') &&
			!(char >= '0' && char <= '9') {
			return false
		}
	}
	return true
}
