package tools

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	keyLength   uint32
	saltLength  uint32
	parallelism uint8
}

var defaultParams = params{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	keyLength:   32,
	saltLength:  16,
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func HashPassword(password string) (string, error) {
	salt, err := generateRandomBytes(defaultParams.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, defaultParams.iterations, defaultParams.memory, defaultParams.parallelism, defaultParams.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultParams.memory,
		defaultParams.iterations,
		defaultParams.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func ComparePasswordAndHash(password, encodedHash string) (bool, error) {
	params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Hash the password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	// Compare the hashes in constant time. The password is correct if the hashes match.
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

func decodeHash(hash string) (p *params, salt, key []byte, err error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &params{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}

	key, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}

	p.keyLength = uint32(len(key))
	p.saltLength = uint32(len(salt))

	return p, salt, key, nil
}

func generateRandomBytes(length uint32) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
