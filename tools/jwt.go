package tools

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	defaultJWTSecret = "netshop-ua-djwt"
	defaultExpire    = 24 * time.Hour
)

// Creates a new JWT token with the given values (claims)
func NewJWTToken(claims jwt.MapClaims) (string, error) {
	var (
		jwtSecret string = TryGetEnv("JWT_SECRET", defaultJWTSecret)
		jwtExpire string = TryGetEnv("JWT_EXPIRE", defaultExpire.String())
	)

	expire, err := time.ParseDuration(jwtExpire)
	if err != nil {
		expire = defaultExpire
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	claims["exp"] = time.Now().Add(expire).Unix()
	claims["iat"] = time.Now().Unix()

	return token.SignedString([]byte(jwtSecret))
}

// Parses a JWT token and returns the claims (values) of the token
func ParseJWTToken(tokenString string) (jwt.MapClaims, error) {
	secret := TryGetEnv("JWT_SECRET", defaultJWTSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}
	return token.Claims.(jwt.MapClaims), nil
}
