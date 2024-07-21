package tools

import (
	"netshop/main/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserTokenClaims struct {
	Id       int64  `json:"id"`
	Type     string `json:"type"`
	Username string `json:"username"`
}

// Creates a new JWT token with the given values (claims)
func NewJWTToken(claims jwt.MapClaims) (string, error) {
	expire, err := time.ParseDuration(config.AppConfig.JwtExpire)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	claims["exp"] = time.Now().Add(expire).Unix()
	claims["iat"] = time.Now().Unix()

	return token.SignedString([]byte(config.AppConfig.JwtSecret))
}

// Parses a JWT token and returns the claims (values) of the token
func ParseJWTToken(tokenString string) (*UserTokenClaims, error) {
	secret := config.AppConfig.JwtSecret
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	claims := token.Claims.(jwt.MapClaims)
	return &UserTokenClaims{
		Id:       int64(claims["id"].(float64)),
		Type:     claims["type"].(string),
		Username: claims["username"].(string),
	}, nil
}
