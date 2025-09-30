package handlers

import (
	"net/http"
	"path/filepath"
	"room-visa/internal/model"
)

type UserHandler struct {
	fs model.FormService 
}

func NewUserHandler(fs model.FormService) *UserHandler {
	return &UserHandler{
        fs: fs,
	}
}

func (uh *UserHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", uh.homePage)
	router.HandleFunc("GET /form", uh.formPage)
	router.HandleFunc("POST /form", uh.readForm)
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

func (uh *UserHandler) readForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

    photo, _, err := r.FormFile("photo")
    if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
    } 

    form := &model.Form{
        Name: r.FormValue("name"), 
        Surname: r.FormValue("surname"),
        Sex: r.FormValue("sex"),
        Ethnicity: r.FormValue("ethnicity"),
        Citizenship: r.FormValue("citizenship"),
        Purpose: r.FormValue("purpose"),
        Photo: photo,
    }

    err = uh.fs.SaveForm(form)
    if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
