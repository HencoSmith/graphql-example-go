package source

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

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

// TODO: complete
// func CreateJWT(userID string) {
// 	key := config.JWT.Key
// 	envKey := os.Getenv("JWT_KEY")
// 	if len(envKey) != 0 {
// 		key = envKey
// 	}
// 	// Setup data to be stored in the token
// 	claims := &Claims{
// 		UserID: userID,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(config.JWT.Expiration * time.Hour).Unix(),
// 		},
// 	}
// 	// Create the token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	// Generate the token string
// 	JWT, err := token.SignedString(key)
// }

// func decodeJWT(JWT string) {
// 	key := config.JWT.Key
// 	claims := &Claims{}

// 	token, err := jwt.ParseWithClaims(JWT, claims, func(token *jwt.Token) (interface{}, error) {
// 		return key, nil
// 	})
// }
