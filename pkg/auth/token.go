package auth

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenData contains token details
type TokenData struct {
	IAT  int64
	ID   int64
	Name string
}

// MakeAuthToken signs a token with given data and returns a token string
func MakeAuthToken(id int64, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"data": jwt.MapClaims{
			"id":   id,
			"name": name,
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	// Shorten token string by removing algorithm details
	tokenParts := strings.Split(tokenString, ".")
	tokenString = tokenParts[1] + "." + tokenParts[2]

	return tokenString, nil
}

// ParseToken prepends the algorithm string to a shortened token,
// parses and verifies it, and returns the parsed data.
func ParseToken(shortToken string) (TokenData, error) {
	var td TokenData
	// Only storing abbreviated tokens in the cookie to reduce cookie size
	// so need to prepend algorithm string (HS256) to received token string.
	tokenString := fmt.Sprintf("%s.%s", os.Getenv("JWT_ALGO"), shortToken)
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return td, err
	}
	if !token.Valid {
		return td, errors.New("Invalid token")
	}

	// iat
	td.IAT = int64(claims["iat"].(float64))
	// data
	claimsData := claims["data"].(map[string]interface{})
	td.ID = int64(claimsData["id"].(float64))
	td.Name = claimsData["name"].(string)

	return td, nil
}
