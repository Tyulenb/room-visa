package app

import (
	"database/sql"
	"log"
	"net/http"
	"room-visa/internal/transport/handlers"
)

type App struct {
    Addr string
    DB *sql.DB
}

func NewApp(port string, db *sql.DB) *App {
    return &App{
        Addr: port,
        DB: db,
    }
}

func (a *App) Run() error {
    router := http.NewServeMux()
    userHandler := handlers.NewUserHandler("")
    userHandler.RegisterRoutes(router)
    adminHandler := handlers.NewAdminHandler("")
    adminHandler.RegisterRoutes(router)

    server := &http.Server{
        Addr: a.Addr,
        Handler: router,
    }

    log.Printf("Server is started on port: %s", a.Addr) 
    return server.ListenAndServe()
}
