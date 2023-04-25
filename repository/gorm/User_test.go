package repo_gorm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Ephilipz/fiberAuth/model"
	"gorm.io/gorm"
)

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

func TestCreate_User(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	repo := NewUserGormRepo(db)
	user := testUsers(1)[0]
	id, err := repo.Create(user)
	if err != nil || id == 0 {
		t.Fatalf("User was not created %s", err)
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

func TestGet_User(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	repo := NewUserGormRepo(db)
	user := testUsers(1)[0]
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

func TestDelete_User(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	repo := NewUserGormRepo(db)
	user := testUsers(1)[0]
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Unable to create the test user %s", err.Error())
	}
	if err := repo.Delete(user.ID); err != nil {
		t.Fatalf("Unable to delete the user %s", err.Error())
	}
	count := int64(0)
	if db.Model(&model.User{}).Where("id = ?", user.ID).Count(&count); count > 0 {
		t.Fatalf("User was not deleted")
	}
}

func TestGetByEmail_User(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	repo := NewUserGormRepo(db)
	user := testUsers(1)[0]
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Unable to create the test user %s", err.Error())
	}
	if getUser := repo.GetByEmail(strings.ToUpper(user.Email)); getUser.ID != user.ID {
		t.Fatalf("Unable to get by email")
	}
}

func TestGetAll_User(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
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

func TestUpdate_User(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	repo := NewUserGormRepo(db)
	user := model.User{
		FirstName: "test",
		LastName:  "last",
		Email:     "test@test.com",
		Password:  []byte("password12345678"),
	}
	err := db.Create(&user).Error
	if err != nil {
		t.Fatalf("Unable to create the test user %s", err.Error())
	}

	userUpdated := model.User{
		Model:     gorm.Model{ID: user.ID},
		FirstName: "newTest",
	}
	err = repo.Update(userUpdated)
	if err != nil {
		t.Fatalf("Unable to update the user %s", err.Error())
	}
	err = db.Model(&user).First(&user).Error
	if err != nil {
		t.Fatalf("Unable to retrieve the user %s", err.Error())
	}
	if user.FirstName != userUpdated.FirstName {
		t.Errorf("updated firstname not matching")
	}
}

func TestGetRoles_User(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
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

func TestUpdateRoles_User(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

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
	if err := db.Create(&user).Association("Roles").Append(&roles[0]); err != nil {
		t.Fatalf("Unable to create the test user %s", err.Error())
	}

	if err := repo.UpdateRoles(user.ID, []uint{roles[1].ID}); err != nil {
		t.Fatalf("Unable to update the user's roles %s", err.Error())
	}

	rolesAfterUpdate := []model.Role{}
	db.Model(&user).Association("Roles").Find(&rolesAfterUpdate)
	if rolesAfterUpdate[0].ID != roles[1].ID {
		t.Fatalf("Roles were not updated correctly")
	}
}
