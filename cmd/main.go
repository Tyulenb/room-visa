package main

import (
	"log"
	"path/filepath"
	"room-visa/internal/app"

	"github.com/joho/godotenv"
)

func main() {

    //Loading enviroment variables
    envPath, err := filepath.Abs("config/.env")
    if err != nil {
        log.Fatalf("Can't load env vars, err: %v ", err)
    }
    if err := godotenv.Load(envPath); err != nil {
        log.Fatalf("Can't load env vars, err: %v ", err)
    }


    app := app.NewApp(":3456", nil)
    app.Run()
}
