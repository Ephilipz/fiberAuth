package route

import (
	"github.com/Ephilipz/fiberAuth/service"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// Returns a middlware handler that verifies the jwt
func ProtectedMiddlware(jwtService service.JWT) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningMethod: "RS256",
		SigningKey:    jwtService.PublicKey(),
	})
}
