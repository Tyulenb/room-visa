package handlers

import (
	"net/http"
	"path/filepath"
)

type UserHandler struct {
    us any
}

func NewUserHandler(us any) *UserHandler {
    return &UserHandler{
        us: us,
    }
}

func (uh *UserHandler) RegisterRoutes(router *http.ServeMux) {
    router.HandleFunc("GET /", uh.homePage)
    router.HandleFunc("GET /form", uh.formPage)
    router.HandleFunc("GET /status", uh.statusPage)
}

func (uh *UserHandler) homePage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    path, err := filepath.Abs("web/index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    http.ServeFile(w, r, path)
}

func (uh *UserHandler) formPage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    path, err := filepath.Abs("web/form.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    http.ServeFile(w, r, path)
}

func (uh *UserHandler) statusPage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    path, err := filepath.Abs("web/request.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.ServeFile(w, r, path)
}

