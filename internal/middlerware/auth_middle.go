package middlerware

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
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
    if token.Valid {
        return nil
    }else {
        return err
    }
}
