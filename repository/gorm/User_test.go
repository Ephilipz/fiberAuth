package repo_gorm

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Ephilipz/fiberAuth/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Unable to initialize db %s", err)
	}
	db.AutoMigrate(&model.User{})
	return db
}

func tearDownDB(t *testing.T) {
	if err := os.Remove("gorm.db"); err != nil {
		t.Fatalf("Unable to delete test db %s", err)
	}
}

func testUsers(count uint8) []model.User {
	var users []model.User
	for i := 0; i < int(count); i++ {
		users = append(users, model.User{
			FirstName: fmt.Sprintf("First%d", i),
			LastName:  fmt.Sprintf("Last%d", i),
			Email:     fmt.Sprintf("test%d@gmail.com", i),
			Password:  []byte("TestPass"),
		})
	}
	return users
}

func TestCreate(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t)
	repo := NewUserGormRepo(db)
	user := testUsers(1)[0]
	id, err := repo.Create(user)
	if id == 0 {
		t.Fatalf("User was not created")
	}
	dbUser := model.User{}
	err = db.Find(&dbUser, id).Error
	if err != nil {
		t.Fatalf("Unable to retrieve the user after create %s", err.Error())
	}
	if dbUser.FirstName != user.FirstName ||
		dbUser.LastName != user.LastName ||
		dbUser.Email != user.Email {
		t.Fatalf("Incorrect fields inserted. Expected %+v got %+v", user, dbUser)
	}

	if _, err = repo.Create(user); err == nil {
		t.Fatalf("Repository inserted two users with the same email")
	}
}

func TestGet(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t)
	repo := NewUserGormRepo(db)
	user := model.User{
		FirstName: "First",
		LastName:  "Last",
		Email:     "test@testing.com",
		Password:  []byte("TestPass"),
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Unable to create the test user %s", err.Error())
	}

	getUser, err := repo.Get(user.ID)
	if err != nil {
		t.Fatalf("Error getting user %s", err)
	}
	if getUser.ID != user.ID {
		t.Fatalf("Mismatching IDs for get")
	}
}

func TestGetByEmail(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t)
	repo := NewUserGormRepo(db)
	user := testUsers(1)[0]
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Unable to create the test user %s", err.Error())
	}
	if getUser, _ := repo.GetByEmail(strings.ToUpper(user.Email)); getUser.ID != user.ID {
		t.Fatalf("Unable to get by email")
	}
}

func TestGetAll(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t)
	repo := NewUserGormRepo(db)
	users := testUsers(20)
	if err := db.CreateInBatches(&users, len(users)).Error; err != nil {
		t.Fatalf("Unable to create users %s", err.Error())
	}
	if getUsers, _ := repo.GetAll(); len(getUsers) != len(users) {
		t.Fatalf("Missing users in get all. Expected %d users; got %d users",
			len(users),
			len(getUsers))
	}
}

func TestGetRoles(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t)
	repo := NewUserGormRepo(db)
	roles := []model.Role{
		{
			Name:      "testRole",
			IsDefault: true,
		},
		{
			Name: "casseRole",
		},
	}
	user := testUsers(1)[0]
	if err := db.CreateInBatches(&roles, len(roles)).Error; err != nil {
		t.Fatalf("Unable to create the test roles %s", err.Error())
	}
	if err := db.Create(&user).Association("Roles").Append(&roles); err != nil {
		t.Fatalf("Unable to create the test user %s", err.Error())
	}

	getRoles, err := repo.GetRoles(user.ID)
	if err != nil {
		t.Fatalf("Error getting roles %s", err.Error())
	}
	if len(getRoles) != len(roles) {
		t.Fatalf("Mismatching roles")
	}
}
