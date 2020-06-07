package actions

import (
	"os"

	"github.com/gofiber/fiber"
)

var (
	admin    = os.Getenv("ADMIN_USER")
	password = os.Getenv("ADMIN_PASSWORD")
)

// Login for user auth
func Login(c *fiber.Ctx) {
	Render(c, "Login", nil, "Login")
}

// Auth the user
func Auth(c *fiber.Ctx) {
	store := sess.Get(c)
	user := c.FormValue("user")
	psd := c.FormValue("password")
	if user == admin && psd == password {
		store.Set("islogined", true)
		store.Save()
		c.Redirect("/admin/pages")
		return
	}
	c.Redirect("/")
}

// Logout out
func Logout(c *fiber.Ctx) {
	store := sess.Get(c)
	store.Destroy()
	c.Redirect("/")
}
