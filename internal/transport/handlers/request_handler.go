package handlers

import (
	"net/http"
	"path/filepath"
)

type RequestHandler struct {
    rs any
}

func NewRequestHandler(rs any) *RequestHandler {
    return &RequestHandler{
        rs: rs,
    }
}

func (rh *RequestHandler) RegisterRoutes(router *http.ServeMux) {
    router.HandleFunc("GET /status", rh.statusPage)
}

func (rh *RequestHandler) statusPage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    path, err := filepath.Abs("web/request.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.ServeFile(w, r, path)
}
