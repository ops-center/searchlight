package auth

import (
	"crypto/sha512"
	"crypto/subtle"
)

// SecureCompare performs a constant time compare of two strings to limit timing attacks.
func SecureCompare(given string, actual string) bool {
	givenSha := sha512.Sum512([]byte(given))
	actualSha := sha512.Sum512([]byte(actual))

	return subtle.ConstantTimeCompare(givenSha[:], actualSha[:]) == 1
}
