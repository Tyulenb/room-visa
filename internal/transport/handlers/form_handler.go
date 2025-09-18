package handlers

import (
	"fmt"
	//	"io"
	"net/http"
)

type FormHandler struct {
}

func NewFormHandler() *FormHandler {
	return &FormHandler{}
}

func (f *FormHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /form", f.readForm)
}

func (f *FormHandler) readForm(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Body is nil", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileHeader := r.MultipartForm.File["photo"]
	file, err := fileHeader[0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(r.MultipartForm.Value)
	fmt.Println(r.MultipartForm.File)
	fmt.Println(file)
}
