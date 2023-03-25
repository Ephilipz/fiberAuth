package repo_gorm

import (
	"os"
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

func TestCreate(t *testing.T) {
	db := setupDB(t)
	defer tearDownDB(t)
	repo := NewUserGormRepo(db)
	user := model.User{
		FirstName: "First",
		LastName:  "Last",
		Email:     "test@testing.com",
		Password:  []byte("TestPass"),
	}
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
