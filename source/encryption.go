package source

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Claims associated with the JWT data stored
type Claims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

// Hash - Hash the text using bcrypt, return the hash or an error
func Hash(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ValidHash - Compare the hash and text using bcrypt, returns true if valid or false along
// with an error if invalid
func ValidHash(text string, hashedText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(text))
	if err != nil {
		return false, err
	}
	return true, nil
}

// CreateJWT return a JWT token for the given user ID input
func CreateJWT(userID string) (string, error) {
	// Read configuration file
	config := GetConfig(".")

	key := config.JWT.Key
	envKey := os.Getenv("JWT_KEY")
	if len(envKey) != 0 {
		key = envKey
	}
	// Setup data to be stored in the token
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.JWT.Expiration * time.Hour).Unix(),
		},
	}
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate the token string

	signingKey := []byte(key)
	return token.SignedString(signingKey)
}

// DecodeJWT decode the specified token and return the associated user ID
func DecodeJWT(JWT string) (string, error) {
	// Read configuration file
	config := GetConfig(".")

	key := config.JWT.Key

	token, err := jwt.ParseWithClaims(JWT, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected token signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", errors.New("Invalid token")
}
