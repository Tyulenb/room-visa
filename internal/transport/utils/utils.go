package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func generateToken(mc jwt.MapClaims, secret string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func GenerateAuthToken() (string, error) {
    admClaims := jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	}
	secret := os.Getenv("ADMIN_AUTH_KEY")

	return generateToken(admClaims, secret)
}

func GenerateVisaToken(req uuid.UUID) (string, error) {
    visaClaims := jwt.MapClaims{
        "sub": "visa",
        "exp": time.Now().Add(time.Hour*32).Unix(),
        "iat": time.Now(),
        "visaReq": req,
    }
    secret := os.Getenv("VISA_KEY")
    return generateToken(visaClaims, secret)
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("The body is empty")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
