package utils

import (
	"strings"
	"unicode"
)



func IsStrongPassword(password string) bool {
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		case strings.ContainsRune("@$!%*?&", ch):
			hasSpecial = true
		}
	}

	return len(password) >= 8 && hasUpper && hasLower && hasNumber && hasSpecial
}
