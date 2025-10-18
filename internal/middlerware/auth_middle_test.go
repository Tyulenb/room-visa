package middlerware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"time"

    "room-visa/internal/crypto"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(t *testing.T, sub string, exp time.Time) (string, error) {
    t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
	})

	secret := os.Getenv("ADMIN_AUTH_KEY")
    token.Header["kid"] = "AuthAdmin"
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

        err := crypto.ValidToken(token)
        if err != nil && !errors.Is(err, tt.expectedErr) {
            t.Fatalf("Unexpected error in test %v\nexpected:%v but got:%v", tt.name, tt.expectedErr, err)
        }
        if err == nil && tt.expectedErr != nil {
            t.Fatalf("Unexpected error in test %v\nexpected:%v but got:%v", tt.name, tt.expectedErr, err)
        }
    }
}

func TestParseTokenFromCookie(t *testing.T) {
    tokenExp, _ := GenerateToken(t, "admin", time.Now().Add(time.Hour))
    tests := []struct {
        name string
        addCookie bool
        cookieName string
        cookieVal string
        expectedErr error
    }{
        {
            name: "Valid",
            addCookie: true,
            cookieName: "AuthAdminToken",
            cookieVal: tokenExp,
            expectedErr: nil,
        },
        {
            name: "Invalid Cookie name",
            addCookie: true,
            cookieName: "RandomName",
            cookieVal: tokenExp,
            expectedErr: http.ErrNoCookie,
        },
        {
            name: "Empty Cookie name",
            addCookie: true,
            cookieName: "",
            cookieVal: tokenExp,
            expectedErr: http.ErrNoCookie,
        },
        {
            name: "No Cookie",
            addCookie: false,
            cookieName: "",
            cookieVal: "",
            expectedErr: http.ErrNoCookie,
        },
    }

    for _, tt := range tests {
        req := httptest.NewRequest(http.MethodGet, "/", nil)
        if tt.addCookie {
            req.AddCookie(&http.Cookie{Name: tt.cookieName, Value: tt.cookieVal})
        }

        token, err := ParseTokenFromCookie(req) 

        if err != nil && !errors.Is(err, tt.expectedErr) {
            t.Fatalf("Unexpected error in test %v\nexpected:%v but got:%v", tt.name, tt.expectedErr, err)
        }
        if err == nil && tt.expectedErr != nil {
            t.Fatalf("Unexpected error in test %v\nexpected:%v but got:%v", tt.name, tt.expectedErr, err)
        }

        if tt.expectedErr == nil && token != tt.cookieVal {
            t.Fatalf("Test:%v\nWrong cookie value\nexpected cookie value: %v but got: %v", tt.name, tt.cookieVal, token)
        }
    }

}

