package repo_gorm

import (
	"testing"

	"github.com/Ephilipz/fiberAuth/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Unable to initialize db %s", err)
	}
	db.AutoMigrate(&model.User{})
	return db
}
