package middlerware

import (
	"errors"
	"os"
	"testing"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

//TO DO
func GenerateToken(t *testing.T, sub string, exp time.Time) (string, error) {
    t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
	})

	secret := os.Getenv("ADMIN_AUTH_KEY")
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}

func TestValidToken(t *testing.T) {
    tests := []struct{
        name string
        sub string
        exp time.Time
        expectedErr error
    }{
        {
            name: "Valid token",
            sub: "admin",
            exp: time.Now().Add(time.Hour),
            expectedErr: nil,
        },
        {
            name: "Expired",
            sub: "admin",
            exp: time.Now().Add(-1*time.Hour),
            expectedErr: jwt.ErrTokenExpired,
        },
    }

    for _, tt := range tests {
        token, _ := GenerateToken(t, tt.sub, tt.exp)

        err := ValidToken(token)
        if err != nil && !errors.Is(err, tt.expectedErr) {
            t.Errorf("Unexpected error in test %v\n expected:%v but got:%v", tt.name, tt.expectedErr, err)
        }
        if err == nil && tt.expectedErr != nil {
            t.Errorf("Unexpected error in test %v\n expected:%v but got:%v", tt.name, tt.expectedErr, err)
        }
    }
}
