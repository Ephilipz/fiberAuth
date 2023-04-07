package database

import (
	"fmt"

	"github.com/Ephilipz/fiberAuth/config"
	"github.com/Ephilipz/fiberAuth/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DB) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn(&cfg)), &gorm.Config{})

	if cfg.ENABLE_MIGRATIONS {
		db.AutoMigrate(&model.User{})
	}

	return db, err
}

func dsn(cfg *config.DB) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.HOST,
		cfg.USER,
		cfg.PASS,
		cfg.NAME,
		cfg.PORT,
		cfg.SSLMODE)
}
