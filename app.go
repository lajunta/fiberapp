package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/template/html"
	"github.com/lajunta/fiberapp/actions"
	"github.com/lajunta/fiberapp/config"
)

var (
	port = os.Getenv("PORT")
)

func initDatabase() {
	config.DBConnect()
	actions.PageCollection = config.DB.Collection("page")
	actions.PageTags = config.Tags
}

func main() {
	initDatabase()
	defer config.DBDisconnect()

	app := fiber.New()

	config.SetupRoutes(app)

	app.Static("/static", "./public", fiber.Static{
		Compress: true,
	})

	app.Settings.ETag = true
	app.Settings.DisableStartupMessage = true
	app.Settings.Templates = html.New("./views/layout/", ".html")

	log.Println("main app is running on", port)
	app.Listen(port)

}
