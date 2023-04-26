package service

import (
	"github.com/Ephilipz/fiberAuth/model"
)

func InitRoles(roleService Role) error {
	if role := roleService.Get(1); role.ID == 1 {
		return nil
	}
	roles := []model.Role{
		{ID: 1, Name: "Superadmin"},
		{ID: 2, Name: "Staff"},
		{ID: 3, Name: "Instructor"},
		{ID: 4, Name: "Student", IsDefault: true},
	}
	return roleService.CreateMultiple(roles)
}

func InitUsers(userService User) error {
	if user := userService.GetByEmail("admin@admin.com"); user.ID != 0 {
		return nil
	}

	insertUser := model.CreateUserDTO{
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  "pass12345678",
		Email:     "admin@admin.com",
	}
	id, err := userService.Create(insertUser)
	if err == nil {
		err = userService.UpdateRoles(id, []uint{1})
	}
	return err
}
