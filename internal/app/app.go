package app

import (
	"database/sql"
	"log"
	"net/http"
	"room-visa/internal/database"
	"room-visa/internal/middlerware"
	"room-visa/internal/service"
	"room-visa/internal/storage"
	"room-visa/internal/transport/handlers"
)

type App struct {
	Addr string
	DB   *sql.DB
    Storage storage.Storage
}

func NewApp(port string, db *sql.DB, store storage.Storage) *App {
	return &App{
		Addr: port,
		DB:   db,
        Storage: store,
	}
}

func (a *App) Run() error {
	router := http.NewServeMux()
	userHandler := handlers.NewUserHandler("")
	userHandler.RegisterRoutes(router)
    
    formRepository := database.NewFormRepository(a.DB)
    formService := service.NewFormService(formRepository, a.Storage)
	formHandler := handlers.NewFormHandler(formService)
	formHandler.RegisterRoutes(router)

	adminRepository := database.NewAdminRepository(a.DB)
	adminService := service.NewAdminService(adminRepository)
	adminHandler := handlers.NewAdminHandler(adminService, formService)
	adminHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:    a.Addr,
		Handler: middlerware.Logging(router),
	}

	log.Printf("Server is started on port: %s", a.Addr)
	return server.ListenAndServe()
}
