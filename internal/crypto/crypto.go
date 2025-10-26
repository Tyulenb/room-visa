package crypto

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateAuthToken() (string, error) {
    admClaims := jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	}
	secret := os.Getenv("ADMIN_AUTH_KEY")
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, admClaims)
    token.Header["kid"] = "AuthAdmin"

	return token.SignedString([]byte(secret)) 
}

func GenerateVisaToken(req uuid.UUID) (string, error) {
    visaClaims := jwt.MapClaims{
        "sub": "visa",
        "exp": time.Now().Add(time.Hour*32).Unix(),
        "iat": time.Now(),
        "visaReq": req,
    }
    secret := os.Getenv("VISA_KEY")
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, visaClaims)
    token.Header["kid"] = "Visa"

	return token.SignedString([]byte(secret)) 
}

func ValidToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, keyFunc)
    if token.Valid {
        return nil
    }else {
        return err
    }
}

func GetVisaTokenClaimReq(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, keyFunc)
    if !token.Valid {
        return "", err 
    }

    claims, ok := token.Claims.(jwt.MapClaims) 
    if !ok {
        return "", fmt.Errorf("Error during parsing claims")
    }
    sub, ok := claims["sub"].(string)
    if !ok || strings.Compare(sub, "visa") != 0 {
        return "", fmt.Errorf("Unsupported token, sub") 
    }
    visaReq, ok := claims["visaReq"].(string)
    if !ok {
        return "", fmt.Errorf("Unsupported token, visaReq") 
    }
    return visaReq, nil
}

func keyFunc(token *jwt.Token) (any, error) {
    kid, ok := token.Header["kid"].(string)
    if !ok {
        return nil, fmt.Errorf("Token header is not valid")
    }
    
    if strings.Compare(kid, "AuthAdmin") == 0 {
        return []byte(os.Getenv("ADMIN_AUTH_KEY")), nil
    }
    if strings.Compare(kid, "Visa") == 0 {
        return []byte(os.Getenv("VISA_KEY")), nil
    }

    return nil, fmt.Errorf("Token header is not valid")
}
