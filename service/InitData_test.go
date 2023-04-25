package service

import (
	"fmt"
	"testing"

	"github.com/Ephilipz/fiberAuth/model"
	repo_gorm "github.com/Ephilipz/fiberAuth/repository/gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestInitRoles(t *testing.T) {
	t.Parallel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Unable to initialize db %s", err)
	}
	db.AutoMigrate(&model.User{})
	repo := repo_gorm.NewRoleGormRepo(db)

	roleService := NewroleService(repo)
	if err = InitRoles(roleService); err != nil {
		t.Errorf("Unable to init roles %s", err.Error())
	}
	if roles, _ := repo.GetAll(); roles == nil || len(roles) == 0 {
		t.Errorf("Initial roles not inserted")
	}
}

func TestInitUsers(t *testing.T) {
	t.Parallel()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Unable to initialize db %s", err)
	}
	db.AutoMigrate(&model.User{})
	repo := repo_gorm.NewUserGormRepo(db)
	roleRepo := repo_gorm.NewRoleGormRepo(db)
	err = roleRepo.CreateMultiple([]model.Role{
		{Model: gorm.Model{ID: 1}, Name: "Admin"},
		{Model: gorm.Model{ID: 2}, Name: "Student", IsDefault: true}},
	)
	if err != nil {
		t.Fatalf("Unable to init roles %s", err.Error())
	}
	userService := NewUserService(repo)
	if err = InitUsers(userService); err != nil {
		t.Fatalf("Unable to init users %s", err.Error())
	}
	users, _ := repo.GetAll()
	if users == nil || len(users) == 0 {
		t.Errorf("Initial users not inserted")
	}
	for _, user := range users {
		roles, _ := repo.GetRoles(user.ID)
		fmt.Println(roles)
		if roles == nil || len(roles) == 0 {
			continue
		}
		for _, role := range roles {
			if role.ID == 1 {
				return
			}
		}
	}
	t.Error("No user with admin role found")
}
