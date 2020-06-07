package actions

import (
	"time"

	"github.com/gofiber/fiber"
	"github.com/gofiber/session"
)

var (
	sessConfig = session.Config{
		Expiration: 30 * time.Minute,
		GCInterval: 30 * time.Minute,
	}
	sess = session.New(sessConfig)
)

// IsAuthed is a middleware check login status
func IsAuthed(c *fiber.Ctx) {
	store := sess.Get(c)
	islogined := store.Get("islogined")
	if islogined == nil {
		c.Redirect("/")
	} else {
		c.Next()
	}
}
