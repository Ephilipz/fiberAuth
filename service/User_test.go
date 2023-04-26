package service

import (
	"testing"

	"github.com/Ephilipz/fiberAuth/model"
	repo_gorm "github.com/Ephilipz/fiberAuth/repository/gorm"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getTestUser() model.User {
	return model.User{
		ID:        1,
		Email:     "test@test.com",
		FirstName: "fist",
		LastName:  "last",
		Password:  []byte("password12345678"),
	}
}

func getTestUserRepo(t *testing.T) userRepository {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Unable to initialize db %s", err)
	}
	db.AutoMigrate(&model.User{})
	return repo_gorm.NewUserGormRepo(db)
}

func TestCreate(t *testing.T) {
	t.Parallel()
	dto := model.CreateUserDTO{
		Email:     "",
		FirstName: "test",
		LastName:  "test",
		Password:  "short",
	}
	repo := getTestUserRepo(t)
	service := NewUserService(repo)

	// missing email
	id, err := service.Create(dto)
	if id != 0 || err == nil {
		t.Error("Expected error for missing email")
	}

	// missing first name
	dto.Email = "test@test.com"
	dto.FirstName = ""
	id, err = service.Create(dto)
	if id != 0 || err == nil {
		t.Error("Expected error for missing first name")
	}

	// missing last name
	dto.FirstName = "test"
	dto.LastName = ""
	id, err = service.Create(dto)
	if id != 0 || err == nil {
		t.Error("Expected error for missing last name")
	}

	// short password
	dto.LastName = "test"
	id, err = service.Create(dto)
	if id != 0 || err == nil {
		t.Error("Expected error for short password")
	}

	// valid
	dto.Password = "test12345678"
	id, err = service.Create(dto)
	if err != nil {
		t.Errorf("Error creating user %s", err.Error())
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	repo := getTestUserRepo(t)
	service := NewUserService(repo)

	id, err := repo.Create(getTestUser())
	if err != nil {
		t.Fatalf("Unable to create user %s", err.Error())
	}

	dto := model.UpdateUserDTO{
		FirstName: "new",
	}
	// missing ID
	err = service.Update(dto)
	if err == nil {
		t.Error("Expected error for missing ID")
	}

	// valid
	dto.ID = id
	err = service.Update(dto)
	if err != nil {
		t.Error("Error updating user &s", err.Error())
	}

	updatedUser, err := repo.Get(id)
	if err != nil {
		t.Errorf("Unable to retrieve user %s", err.Error())
	}
	if updatedUser.FirstName != dto.FirstName {
		t.Error("Error first name not updated")
	}
}

func TestValidateCredentials(t *testing.T) {
	t.Parallel()
	repo := getTestUserRepo(t)
	service := NewUserService(repo)

	user := getTestUser()
	originalPassword := string(user.Password)
	user.Password, _ = bcrypt.GenerateFromPassword(user.Password, 2)
	id, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Unable to create user %s", err.Error())
	}

	// incorrect email
	if isValid, _ := service.ValidateCredentials("nonExistantEmail", "somePassword"); isValid {
		t.Error("Expected invalid credentials for non existant email")
	}

	// incorrect password
	if isValid, _ := service.ValidateCredentials(user.Email, "wrongPassword"); isValid {
		t.Error("Expected invalid credentials for wrong password")
	}

	// valid
	if isValid, resultId := service.ValidateCredentials(user.Email, originalPassword); !isValid || resultId != id {
		t.Errorf("Expected valid credentials and correct id but got valid: %v, id: %d", isValid, resultId)
	}
}
