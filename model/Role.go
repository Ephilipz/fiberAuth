package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name      string `gorm:"unique;not null;default:null"`
	IsDefault bool   `gorm:"not null;default:false"`
}
