package handler

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Ephilipz/fiberAuth/model"
	"github.com/Ephilipz/fiberAuth/service"
	"github.com/gofiber/fiber/v2"
)

type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(userService service.User, jwtService service.JWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		loginDTO := loginDTO{}
		if err := c.BodyParser(&loginDTO); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "invalid login object"})
		}

		if err := validateLoginDTO(loginDTO); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": err.Error()})
		}

		valid, userId := userService.ValidateCredentials(loginDTO.Email, loginDTO.Password)

		if !valid {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "incorrect username or password"})
		}

		return sendTokens(c, userService, jwtService, userId)
	}
}

func validateLoginDTO(userDTO loginDTO) error {
	err := []error{}
	if len(userDTO.Email) == 0 {
		err = append(err, fmt.Errorf("Email cannot be empty"))
	}
	if len(userDTO.Password) == 0 {
		err = append(err, fmt.Errorf("Password cannot be empty"))
	}

	if len(err) == 0 {
		return nil
	}

	return errors.Join(err...)
}

func Register(userService service.User, jwtService service.JWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userDTO := model.CreateUserDTO{}
		if err := c.BodyParser(&userDTO); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "invalid registration object"})
		}

		if err := validateRegistrationDTO(userDTO); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": err.Error()})
		}

		userId, err := userService.Create(model.CreateUserDTO(userDTO))
		if err != nil {
			return c.Status(fiber.ErrConflict.Code).
				JSON(fiber.Map{"message": "unable to register. This user likely exists"})
		}

		return sendTokens(c, userService, jwtService, userId)
	}
}

func validateRegistrationDTO(userDTO model.CreateUserDTO) error {
	err := []error{}
	if len(userDTO.Email) == 0 {
		err = append(err, fmt.Errorf("Email cannot be empty"))
	}
	if len(userDTO.FirstName) == 0 {
		err = append(err, fmt.Errorf("FirstName cannot be empty"))
	}
	if len(userDTO.LastName) == 0 {
		err = append(err, fmt.Errorf("Lastname cannot be empty"))
	}
	if len(userDTO.Password) < 8 {
		err = append(err, fmt.Errorf("Password cannot be less than 8 characters"))
	}

	if len(err) == 0 {
		return nil
	}
	return errors.Join(err...)
}

func Refresh(userService service.User, jwtService service.JWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.FormValue("token")
		if len(token) == 0 {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "Token cannot be empty"})
		}
		claims, err := jwtService.ClaimsFromJWT(token)
		if err != nil || len(claims) == 0 {
			return c.Status(fiber.ErrBadRequest.Code).
				JSON(fiber.Map{"message": "Invalid token"})
		}

		userId := uint(claims["sub"].(float64))
		return sendTokens(c, userService, jwtService, userId)
	}
}

func sendTokens(c *fiber.Ctx, userService service.User, jwtService service.JWT, userId uint) error {
	if userId == 0 {
		return c.Status(fiber.ErrNotAcceptable.Code).
			JSON(fiber.Map{"message": "Unable to find user"})
	}

	roles, err := userService.GetRoles(userId)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).
			JSON(fiber.Map{"message": "Unable to get roles"})
	}

	roleNamesStr := strings.Join(roleNames(roles), ",")
	claims := map[string]any{
		"sub":   userId,
		"roles": roleNamesStr,
	}

	var access, refresh string
	access, err = jwtService.GenerateJWT(claims, time.Minute*10)
	refresh, err = jwtService.GenerateJWT(claims, time.Hour*24)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).
			JSON(fiber.Map{"message": "Unable to generate JWT"})
	}

	return c.JSON(fiber.Map{"roles": roleNamesStr, "access": access, "refresh": refresh})
}

func roleNames(roles []model.Role) []string {
	roleNames := []string{}
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames
}
