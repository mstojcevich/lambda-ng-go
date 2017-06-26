package user

import (
	"crypto/sha256"
	"errors"
	"strconv"

	"encoding/base64"

	"fmt"

	"strings"

	"github.com/dchest/uniuri"
	"golang.org/x/crypto/pbkdf2"
)

// Handles password hashing / checking
// Right now just ensures compatibility w/ existing Lambda passwords
// so supports Python's passlib's pbkdf2_sha256 hashing
// which is secure for now.
// Because of its limited compatibility, it's not a good idea to use this for your own project.

// HashPassword hashes a password with a secure hashing algorithm
// The produced string is compatible with the PHC string format: https://github.com/P-H-C/phc-string-format/blob/master/phc-sf-spec.md
// The specific algorithm used should not be relied on, though right now it is pbkdf2_sha256.
func HashPassword(password string) string {
	salt := genSalt(16) // Length 16 salt for now

	return hashPassword(password, salt, 11949, 32)
}

// CheckPassword checks if a plaintext (raw) password matches a hashed password
// It supports hashed passwords in the PHC string format: https://github.com/P-H-C/phc-string-format/blob/master/phc-sf-spec.md
func CheckPassword(rawPassword string, hashedPassword string) (bool, error) {
	dollarSplit := strings.Split(hashedPassword, "$")

	// From django docs: <algorithm>$<iterations>$<salt>$<hash>

	if len(dollarSplit) < 4 {
		return false, errors.New("Invalid password format. Not enough parameters")
	}

	iterations, err := strconv.Atoi(dollarSplit[1])
	if err != nil {
		return false, err
	}

	algorithm := dollarSplit[0]
	if algorithm != "pbkdf2_sha256" { // For now we only support this one algorithm
		return false, errors.New("Unsupported algorithm: " + algorithm)
	}

	salt := dollarSplit[2]
	pw := dollarSplit[3]

	// Need to decode the password to get the length to pass to pbkdf2
	// probably can do math to figure it out instead of the full decode process
	// TODO
	decodedPw, err := base64.StdEncoding.DecodeString(pw)
	if err != nil {
		return false, err
	}

	size := len(decodedPw)

	hashedInput := pbkdf2.Key([]byte(rawPassword), []byte(salt), iterations, size, sha256.New)

	return base64.StdEncoding.EncodeToString(hashedInput) == pw, nil
}

// hashPassword hashes a password with pbkdf2
func hashPassword(password string, salt string, iterations int, length int) string {
	// Right now just emulate what django and passlib do

	// Hash the password with pbdf2
	encPass := pbkdf2.Key([]byte(password), []byte(salt), iterations, length, sha256.New)
	// Base64 encode the hashed password as that's what the format expects
	hashStr := base64.StdEncoding.EncodeToString(encPass)

	// From django docs: <algorithm>$<iterations>$<salt>$<hash>
	return fmt.Sprintf("%s$%s$%s$%s", "pbkdf2_sha256", strconv.Itoa(iterations), salt, hashStr)
}

// Generates a pseudorandom alphanumeric salt of the specified length
func genSalt(length int) string {
	return uniuri.NewLen(length)
}
