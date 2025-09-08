package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)
func GenerateToken() (string, error) {
    // Create a new token object, specifying signing method and the claims
    // you would like it to contain.
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"sub": "admin",
	"exp": time.Now().Add(time.Hour*1).Unix(),
    "iat": time.Now().Unix(),
    })
    // Sign and get the complete encoded token as a string using the secret
    tokenString, err := token.SignedString([]byte("secret"))

    return tokenString, err
}

func keyFunc(token *jwt.Token) (any, error) {
    return []byte(os.Getenv("SECRET_KEY")), nil
}

func ValidToken(tokenString string) error {
    token, err := jwt.Parse(tokenString, keyFunc)
    if err != nil {
        return err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        if claims["exp"].(float64) < float64(time.Now().Unix())  {
            return fmt.Errorf("Token expired") 
        }
	    return nil 
    } else {
        return fmt.Errorf("Incorrect claims")
    }
}

