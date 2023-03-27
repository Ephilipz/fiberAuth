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

func testGet(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(t)
	repo := NewRoleGormRepo(db)
	role := testRoles(1)[0]
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("Unable to create role to test %s", err.Error())
	}
	if getRole := repo.Get(role.ID); getRole.Name != role.Name {
		t.Fatalf("Mismatching roles")
	}
}
