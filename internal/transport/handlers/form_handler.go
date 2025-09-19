package handlers

import (
	"net/http"
	"room-visa/internal/model"
)

type FormHandler struct {
    fs model.FormService
}

func NewFormHandler(fs model.FormService) *FormHandler {
	return &FormHandler{
        fs: fs,
    }
}

func (f *FormHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /form", f.readForm)
}

func (f *FormHandler) readForm(w http.ResponseWriter, r *http.Request) {
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

    err = f.fs.SaveForm(form)
    if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
