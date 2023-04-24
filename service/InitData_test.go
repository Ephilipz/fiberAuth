package service

import (
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

	userService := NewUserService(repo)
	if err = InitUsers(userService); err != nil {
		t.Errorf("Unable to init users %s", err.Error())
	}
	if users, _ := repo.GetAll(); users == nil || len(users) == 0 {
		t.Errorf("Initial users not inserted")
	}
}
