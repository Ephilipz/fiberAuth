package route

import (
	"github.com/Ephilipz/fiberAuth/handler"
	"github.com/Ephilipz/fiberAuth/service"
	"github.com/gofiber/fiber/v2"
)

func Auth(app fiber.Router, userService service.User, jwtService service.JWT) {
	app.Post("/login", handler.Login(userService, jwtService))
	app.Post("/register", handler.Register(userService, jwtService))
	// app.Get("/loginapple", handler.LoginApple(userService))
	// app.Get("/logingoogle", handler.LoginGoogle(userService))
}
