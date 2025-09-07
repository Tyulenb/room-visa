package handlers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

type AdminHandler struct {
    us any
}

func NewAdminHandler(us any) *AdminHandler {
    return &AdminHandler{
        us: us,
    }
}

func (a *AdminHandler) RegisterRoutes(router *http.ServeMux) {
    router.HandleFunc("GET /auth",a.adminAuthPage)
    router.HandleFunc("POST /auth",a.authAdmin)
}

//Sends admin auth page
func (a *AdminHandler) adminAuthPage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    path, err := filepath.Abs("web/admin_auth.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    http.ServeFile(w, r, path)
}

func (a *AdminHandler) authAdmin(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    defer r.Body.Close()

    //Form in body retruns in format 'login=...&password=...' 
    //We need to split it to get login and password
    bodySplit := strings.Split(string(body), "&")
    login := strings.Split(bodySplit[0], "=")[1]
    password := strings.Split(bodySplit[1], "=")[1]
    fmt.Fprintf(w, "Login: %s, Password: %s", login, password)
}
