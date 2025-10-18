package middlerware

import (
	"log"
	"net/http"
    "room-visa/internal/crypto"
)

func AuthAdminMiddle(next http.Handler ) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, err := ParseTokenFromCookie(r)
        if err != nil {
            log.Println("Error: ", err)
            http.Error(w, "Access denied", http.StatusUnauthorized)
            return 
        }
        
        err = crypto.ValidToken(token)
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

