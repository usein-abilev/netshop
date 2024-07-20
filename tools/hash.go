package tools

import "fmt"

import (
	"strconv"

	"golang.org/x/crypto/argon2"
)

const (
	DefaultArgon2Salt    = "netshop-ua-dsalt"
	DefaultArgon2Threads = 4
)

// Hashes the password using the Argon2 algorithm
func HashPassword(password string) (string, error) {
	salt := TryGetEnv("ARGON2_SALT", DefaultArgon2Salt)
	threads, err := strconv.Atoi(TryGetEnv("ARGON2_THREADS", fmt.Sprint(DefaultArgon2Threads)))
	if err != nil {
		threads = DefaultArgon2Threads
	}
	hash := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, uint8(threads), 32)
	return string(hash), nil
}
