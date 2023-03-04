package route

import (
	"github.com/Ephilipz/fiberAuth/handler"
	"github.com/Ephilipz/fiberAuth/service"
	"github.com/gofiber/fiber/v2"
)

func User(app fiber.Router, userService service.User) {
	app.Get("/", handler.GetUsers(userService))
	app.Get("/:id", handler.GetUser(userService))
	app.Get("/:id/roles", handler.GetUserRoles(userService))
	app.Post("/updateroles", handler.UpdateUserRoles(userService))
}
