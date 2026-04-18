package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ParseAccessToken(accessToken string) (int64, int64, error) {
	parser := &jwt.Parser{}
	token, _, err := parser.ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return 0, 0, err
	}
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, fmt.Errorf("invalid claims")
	}
	exp, ok := payload["exp"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("exp not found")
	}
	iat, ok := payload["iat"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("iat not found")
	}
	return int64(exp), int64(iat), nil

}
func GenAccessToken(exp int64, iat int64) (string, error) {

	claims := jwt.MapClaims{
		"exp": exp,
		"iat": iat,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte("my_super_secret_test_key")
	signedTokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedTokenString, nil
}
