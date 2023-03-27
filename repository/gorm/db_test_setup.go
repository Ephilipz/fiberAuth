package repo_gorm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Ephilipz/fiberAuth/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dbName := fmt.Sprintf("testGorm%s", t.Name())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("Unable to initialize db %s", err)
	}
	db.AutoMigrate(&model.User{})
	return db
}

func tearDownTestDB(t *testing.T) {
	dbName := fmt.Sprintf("testGorm%s", t.Name())
	if err := os.Remove(dbName); err != nil {
		t.Fatalf("Unable to delete test db %s", err)
	}
}
