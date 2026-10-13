// Package auth handles password hashing, session cookies and CSRF tokens.
//
// Everything here uses only the standard library: crypto/pbkdf2 landed in
// Go 1.24, so no third-party KDF is needed.
package auth

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	// OWASP's current floor for PBKDF2-HMAC-SHA256.
	pbkdf2Iterations = 600_000
	saltBytes        = 16
	keyBytes         = 32
)

var (
	ErrInvalidHash = errors.New("auth: malformed password hash")
	ErrMismatch    = errors.New("auth: password does not match")
)

// HashPassword returns a self-describing hash string:
//
//	pbkdf2-sha256$<iterations>$<salt-b64>$<key-b64>
//
// Storing the parameters alongside the digest means the cost can be raised
// later without invalidating existing accounts.
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("auth: read salt: %w", err)
	}

	key, err := pbkdf2.Key(sha256.New, password, salt, pbkdf2Iterations, keyBytes)
	if err != nil {
		return "", fmt.Errorf("auth: derive key: %w", err)
	}

	return fmt.Sprintf("pbkdf2-sha256$%d$%s$%s",
		pbkdf2Iterations,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key)), nil
}

// VerifyPassword re-derives the key with the stored parameters and compares in
// constant time.
func VerifyPassword(password, encoded string) error {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 || parts[0] != "pbkdf2-sha256" {
		return ErrInvalidHash
	}

	iterations, err := strconv.Atoi(parts[1])
	if err != nil || iterations < 1000 {
		return ErrInvalidHash
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return ErrInvalidHash
	}
	want, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return ErrInvalidHash
	}

	got, err := pbkdf2.Key(sha256.New, password, salt, iterations, len(want))
	if err != nil {
		return ErrInvalidHash
	}
	if subtle.ConstantTimeCompare(got, want) != 1 {
		return ErrMismatch
	}
	return nil
}

// PasswordProblem validates strength and returns a translation key describing
// the first failure, or "" when the password is acceptable.
func PasswordProblem(pw string) string {
	if len([]rune(pw)) < 8 {
		return "err.password_short"
	}
	if len(pw) > 256 {
		return "err.password_long"
	}

	var hasLetter, hasDigit bool
	for _, r := range pw {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return "err.password_weak"
	}
	return ""
}
