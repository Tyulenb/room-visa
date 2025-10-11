package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"room-visa/internal/middlerware"
	"room-visa/internal/model"
	"room-visa/internal/transport/utils"
	"time"

	"github.com/google/uuid"
)

type AdminHandler struct {
	as model.AdminService
    fs model.FormService
}

func NewAdminHandler(as model.AdminService, fs model.FormService) *AdminHandler {
	return &AdminHandler{
		as: as,
        fs: fs,
	}
}

func (a *AdminHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /auth", a.adminAuthPage)
	router.HandleFunc("POST /auth", a.authAdmin)
	router.HandleFunc("POST /addAdmin", a.addAdmin)
    router.Handle("GET /requests", middlerware.AuthAdminMiddle(http.HandlerFunc(a.listRequests)))
    router.Handle("GET /requests/check", middlerware.AuthAdminMiddle(http.HandlerFunc(a.checkRequest)))
}

// Sends admin auth page
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

    cookie := http.Cookie {
        Name: "AuthAdminToken",
        Value: token,
        Expires: time.Now().Add(12*time.Hour),
        Path: "/",
        HttpOnly: true,
    }

    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/requests", http.StatusSeeOther)
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

func (a *AdminHandler) listRequests(w http.ResponseWriter, r *http.Request) {
	data, err := a.fs.GetForms()
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadGateway)
		log.Println("GetForms:", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, `<!doctype html><html><head><meta charset=\"utf-8\"><title>Requests</title><script src="https://unpkg.com/htmx.org@1.9.4"></script></head><body>`)

    tmpl, err := template.ParseFiles("web/request_template.html")
    if err != nil {
        http.Error(w, "Something went wrong", http.StatusBadGateway)
        log.Println("ParseFiles:", err)
        return
    }

	for i, v := range data {
		photo, err := a.fs.LoadFormPhoto(v.Photo) //photo in base64 format
		if err != nil {
            http.Error(w, "Something went wrong", http.StatusBadGateway)
            log.Println("LoadFormPhoto:", err)
			continue
		}
        data[i].Photo = photo //Now represents image encoding in base 64, for template
        tmpl.Execute(w, data[i])
	}

	fmt.Fprintln(w, "</body></html>")
}

func (a *AdminHandler) checkRequest(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()
    status := values.Get("result")
    id := values.Get("id")
    uuid, err := uuid.Parse(id)
    if err != nil {
        http.Error(w, "Something went wrong", http.StatusBadGateway)
        log.Println("uuidParse:", err)
        return
    }
    err = a.fs.ChangeFormStatus(uuid, status)
    if err != nil {
        http.Error(w, "Something went wrong", http.StatusBadGateway)
        log.Println("ChangeFormStatus:", err)
        return
    }
}

