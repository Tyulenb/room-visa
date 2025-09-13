package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"room-visa/internal/model"
	"room-visa/internal/transport/utils"
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
    router.HandleFunc("POST /addAdmin",a.addAdmin)
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
    login, password := r.FormValue("login"), r.FormValue("password")

    err := a.as.CheckPassword(login, password)
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

func (a *AdminHandler) addAdmin(w http.ResponseWriter, r *http.Request) {
    admin := new(model.Admin)
    if err := utils.ParseJSON(r, admin); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    
    adm, err := a.as.CreateAdmin(admin.Login, admin.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
   
    if err := utils.WriteJSON(w, http.StatusOK, adm); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
}
