package main

import "room-visa/internal/app"

func main() {
    app := app.NewApp(":3456", nil)
    app.Run()
}
