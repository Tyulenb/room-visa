package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
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
    fmt.Fprintf(w, "POST")
}
