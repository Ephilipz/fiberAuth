package model

type Role struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null;default:null"`
	IsDefault bool   `gorm:"not null;default:false"`
}
