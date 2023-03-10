package service

import (
	"github.com/Ephilipz/fiberAuth/model"
	"gorm.io/gorm"
)

// TODO : insert initial users and roles from yaml

func InitRoles(roleService Role) error {
	roles := []model.Role{
		{Model: gorm.Model{ID: 1}, Name: "Superadmin"},
		{Model: gorm.Model{ID: 2}, Name: "Staff"},
		{Model: gorm.Model{ID: 3}, Name: "Instructor"},
		{Model: gorm.Model{ID: 4}, Name: "Student", IsDefault: true},
	}

	return roleService.CreateMultiple(roles)
}

func InitUsers(userService User) error {
	user := model.CreateUserDTO{
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  "pass12345678",
		Email:     "admin@admin.com",
	}

	_, err := userService.Create(user)
	return err
}
