package repo_gorm

import (
	"fmt"
	"testing"

	"github.com/Ephilipz/fiberAuth/model"
)

func testRoles(count uint8) []model.Role {
	roles := []model.Role{}
	for i := 0; i < int(count); i++ {
		roles = append(roles, model.Role{
			Name: fmt.Sprintf("TestRole%d", i),
		})
	}
	return roles
}

func TestGet_Role(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	repo := NewRoleGormRepo(db)
	role := testRoles(1)[0]
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("Unable to create role to test %s", err.Error())
	}
	if getRole := repo.Get(role.ID); getRole.Name != role.Name {
		t.Fatalf("Mismatching roles")
	}
}

func TestGetAll_Role(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	repo := NewRoleGormRepo(db)
	roles := testRoles(3)
	if err := db.Create(&roles).Error; err != nil {
		t.Fatalf("Unable to create roles to test %s", err.Error())
	}
	getRoles, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Error retrieving roles: %s", err.Error())
	}
	if len(getRoles) != len(roles) {
		t.Fatalf("Mismatching number of roles retrieved")
	}
	for i, role := range roles {
		if role.ID != getRoles[i].ID {
			t.Fatalf("Expected role ID %d but got %d", role.ID, getRoles[i].ID)
		}
		if role.Name != getRoles[i].Name {
			t.Fatalf("Expected role name %s but got %s", role.Name, getRoles[i].Name)
		}
	}
}

func TestGetByName_Role(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	repo := NewRoleGormRepo(db)
	role := testRoles(1)[0]
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("Unable to create role to test %s", err.Error())
	}
	getRole := repo.GetByName(role.Name)
	if getRole.ID != role.ID {
		t.Fatalf("Mismatching roles")
	}
}

func TestCreate_Role(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	repo := NewRoleGormRepo(db)
	role := testRoles(1)[0]
	if err := repo.Create(role); err != nil {
		t.Fatalf("Error creating role: %s", err.Error())
	}
	getRole := model.Role{}
	if err := db.First(&getRole).Error; err != nil {
		t.Fatalf("Error retrieving created role: %s", err.Error())
	}
	if getRole.Name != role.Name {
		t.Fatalf("Mismatching roles")
	}
}

func TestCreateMultiple_Role(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	repo := NewRoleGormRepo(db)
	roles := testRoles(3)
	if err := repo.CreateMultiple(roles); err != nil {
		t.Fatalf("Error creating roles: %s", err.Error())
	}
	getRoles := []model.Role{}
	if err := db.Find(&getRoles).Error; err != nil {
		t.Fatalf("Error retrieving created roles: %s", err.Error())
	}
	if len(getRoles) != len(roles) {
		t.Fatalf("Mismatching number of roles retrieved")
	}
}

func TestDelete_Role(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	repo := NewRoleGormRepo(db)
	role := testRoles(1)[0]
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("Unable to create role to test %s", err.Error())
	}
	if err := repo.Delete(role.ID); err != nil {
		t.Fatalf("Error deleting role: %s", err.Error())
	}
	getRole := model.Role{}
	if err := db.First(&getRole, role.ID).Error; err == nil {
		t.Fatalf("Role was not deleted")
	}
}

func TestUpdate_Role(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	repo := NewRoleGormRepo(db)
	role := testRoles(1)[0]
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("Unable to create role to test %s", err.Error())
	}
	role.Name = "UpdatedRole"
	if err := repo.Update(role); err != nil {
		t.Fatalf("Error updating role: %s", err.Error())
	}
	getRole := model.Role{}
	if db.First(&getRole, role.ID); getRole.Name != role.Name {
		t.Fatalf("Role name was not updated")
	}
}
