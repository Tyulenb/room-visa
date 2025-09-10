package handlers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"room-visa/internal/model"
	"room-visa/internal/transport/utils"
	"strings"
)

type AdminHandler struct {
    as model.AdminService
}

func NewAdminHandler(as model.AdminService) *AdminHandler {
    return &AdminHandler{
        as: as,
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


    err = a.as.CheckPassword(login, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }

    token, err := utils.GenerateToken()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }

    fmt.Fprintf(w, "Login: %s, Password: %s, Token: %s", login, password, token)
}
