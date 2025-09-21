package middlerware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func AuthAdminMiddle(next http.Handler ) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, err := ParseTokenFromCookie(r)
        if err != nil {
            log.Println("Error: ", err)
            http.Error(w, "Access denied", http.StatusUnauthorized)
            return 
        }
        
        err = ValidToken(token)
        if err != nil {
            log.Println("Error: ", err)
            http.Error(w, "Access denied", http.StatusUnauthorized)
            return 
        }
        next.ServeHTTP(w, r)
    })
}

func ParseTokenFromCookie(r *http.Request) (string, error){
    cookie, err := r.Cookie("AuthAdminToken")
    if err != nil {
        return "", err 
    }
    if err := cookie.Valid(); err != nil {
        return "", err 
    }
    
    return cookie.Value, nil
}

func keyFunc(token *jwt.Token) (any, error) {
	return []byte(os.Getenv("ADMIN_AUTH_KEY")), nil
}

func ValidToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
        exp, ok := claims["exp"].(float64)
        if !ok {
            return fmt.Errorf("invalid exp claim")
        }
		if exp < float64(time.Now().Unix()) {
			return fmt.Errorf("Token expired")
		}
		return nil
	} else {
		return fmt.Errorf("Incorrect claims")
	}
}
