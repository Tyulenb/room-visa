package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"room-visa/internal/app"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	connStr := os.Getenv("DB")
	db, err := sql.Open("postgres", connStr)
	if err := dbPing(db); err != nil {
		log.Fatalf("Database is unavailable, err:%v", err)
	}

	app := app.NewApp(":3456", db)
	app.Run()
}

func dbPing(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	log.Println("DB sucessfully connected")
	return nil
}
