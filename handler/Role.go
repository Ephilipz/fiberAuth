package handler

import (
	"strconv"

	"github.com/Ephilipz/fiberAuth/model"
	"github.com/Ephilipz/fiberAuth/service"
	"github.com/gofiber/fiber/v2"
)

func GetRole(roleService service.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "Invalid Id"})
		}
		role := roleService.Get(uint(id))
		if role.ID == 0 {
			return c.Status(fiber.ErrNotFound.Code).
				JSON(fiber.Map{"message": "No role found"})
		}

		return c.JSON(role)
	}
}

func GetRoles(roleService service.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, err := roleService.GetAll()
		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "Unable to get roles"})
		}

		return c.JSON(roles)
	}
}

func CreateRole(roleService service.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		if len(name) == 0 {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "no name provided"})
		}

		role := model.Role{Name: name}

		err := roleService.Create(role)

		if err != nil {
			return c.Status(fiber.ErrInternalServerError.Code).
				JSON(fiber.Map{"message": "Unable to create role"})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
