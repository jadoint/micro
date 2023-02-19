package token

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Create signs a token with given data and returns a token string
func Create(dataClaim *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":  time.Now().Unix(),
		"data": dataClaim,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Parse prepends the algorithm string to a shortened token, if needed.
// If the token is not shortened, just parse and verify it, and
// return the parsed data.
//
// Example parsing of returned claims:
// td.IAT = int64(claims["iat"].(float64))
// claimsData := claims["data"].(map[string]interface{})
// td.ID = int64(claimsData["id"].(float64))
// td.Name = claimsData["name"].(string)
func Parse(token string) (jwt.MapClaims, error) {
	// Determine if the token is a short token or a full one.
	tokenParts := strings.Split(token, ".")
	isFullToken := len(tokenParts) == 3
	tokenString := token
	if !isFullToken {
		// Only storing abbreviated tokens in the cookie to reduce cookie size
		// so need to prepend algorithm string (HS256) to received token string.
		tokenString = fmt.Sprintf("%s.%s", os.Getenv("JWT_ALGO"), token)
	}
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(parsedToken *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return claims, err
	}
	if !parsedToken.Valid {
		return claims, errors.New("Invalid token")
	}
	return claims, nil
}
