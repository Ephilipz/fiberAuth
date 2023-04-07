package route

import (
	"strings"

	"github.com/Ephilipz/fiberAuth/service"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

// Returns a middlware handler that verifies the jwt
func ProtectedMiddlware(jwtService service.JWT) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningMethod: "RS256",
		SigningKey:    jwtService.PublicKey(),
	})
}

func RestrictedMiddleware(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
		roles := claims["roles"].(string)
		if !strings.Contains(roles, role) {
			return c.Status(fiber.ErrForbidden.Code).JSON(
				fiber.Map{"message": "Invalid role"},
			)
		}
		return c.Next()
	}
}
