package handler

import (
	"github.com/Ephilipz/fiberAuth/service"
	"github.com/gofiber/fiber/v2"
)

func GetUser(userService service.User) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "Invalid id"})
		}
		user, err := userService.Get(uint(id))

		if err != nil {
			return c.Status(fiber.ErrNotFound.Code).
				JSON(fiber.Map{"message": "Unable to get user"})
		}

		return c.JSON(user)
	}
}

func GetUsers(userService service.User) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := userService.GetAll()
		if err != nil {
			return c.Status(fiber.ErrNotFound.Code).
				JSON(fiber.Map{"message": "Unable to get users"})
		}

		return c.JSON(users)
	}
}

type updateUserRolesDTO struct {
	UserId  uint   `json:"userId"`
	RoleIds []uint `json:"roleIds"`
}

func UpdateUserRoles(userService service.User) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := updateUserRolesDTO{}
		if err := c.BodyParser(&dto); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "Invalid user or roles"})
		}

		err := userService.UpdateRoles(dto.UserId, dto.RoleIds)
		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "Unable to update roles"})
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

func GetUserRoles(userService service.User) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "Invalid user id"})
		}
		roles, err := userService.GetRoles(uint(id))
		if err != nil {
			return c.Status(fiber.ErrNotFound.Code).
				JSON(fiber.Map{"message": "Unable to get roles"})
		}

		return c.JSON(roles)
	}
}
