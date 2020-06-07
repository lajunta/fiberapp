package config

import (
	"github.com/gofiber/fiber"
	"github.com/lajunta/fiberapp/actions"
)

func security(c *fiber.Ctx) {
	// Set some security headers:
	c.Set("X-XSS-Protection", "1; mode=block")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Download-Options", "noopen")
	c.Set("Strict-Transport-Security", "max-age=5184000")
	c.Set("X-Frame-Options", "SAMEORIGIN")
	c.Set("X-DNS-Prefetch-Control", "off")

	// Go to next middleware:
	c.Next()
}

// SetupRoutes is
func SetupRoutes(app *fiber.App) {
	app.Use(security)

	app.Get("/", actions.Home)
	app.Get("/pages/show/:id", actions.ShowPage)
	app.Get("/tags/:name", actions.TagList)
	app.Get("/login", actions.Login)
	app.Post("/auth", actions.Auth)
	app.Get("/logout", actions.Logout)

	admin := app.Group("/admin", actions.IsAuthed)
	admin.Get("/pages/new", actions.NewPage)
	admin.Get("/pages", actions.GetPages)
	admin.Get("/pages/edit/:id", actions.EditPage)
	admin.Get("/pages/delete/:id", actions.DeletePage)
	admin.Post("/pages", actions.CreatePage)
	admin.Post("/pages/:id", actions.UpdatePage)
	admin.Get("/search", actions.GetPages)

}
