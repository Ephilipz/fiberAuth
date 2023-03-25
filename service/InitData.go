package service

import (
	"github.com/Ephilipz/fiberAuth/model"
	"gorm.io/gorm"
)

func InitRoles(roleService Role) error {
	if role, _ := roleService.Get(1); role.ID == 1 {
		return nil
	}
	roles := []model.Role{
		{Model: gorm.Model{ID: 1}, Name: "Superadmin"},
		{Model: gorm.Model{ID: 2}, Name: "Staff"},
		{Model: gorm.Model{ID: 3}, Name: "Instructor"},
		{Model: gorm.Model{ID: 4}, Name: "Student", IsDefault: true},
	}
	return roleService.CreateMultiple(roles)
}

func InitUsers(userService User) error {
	if user, _ := userService.GetByEmail("admin@admin.com"); user.ID != 0 {
		return nil
	}
	insertUser := model.CreateUserDTO{
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  "pass12345678",
		Email:     "admin@admin.com",
	}
	_, err := userService.Create(insertUser)
	return err
}
