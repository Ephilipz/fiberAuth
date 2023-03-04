package route

import (
	"github.com/Ephilipz/fiberAuth/handler"
	"github.com/Ephilipz/fiberAuth/service"
	"github.com/gofiber/fiber/v2"
)

func Role(app fiber.Router, roleService service.Role) {
	app.Get("/", handler.GetRoles(roleService))
	app.Get("/:id", handler.GetRole(roleService))
	app.Post("/create", handler.CreateRole(roleService))
}
