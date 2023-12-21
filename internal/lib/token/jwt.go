package token

import "github.com/golang-jwt/jwt/v5"

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func ParseToken(token string, secret string) (string, error) {

	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := parsedToken.Claims.(*tokenClaims)
	if !ok {
		return "", err
	}

	return claims.UserID, nil
}
