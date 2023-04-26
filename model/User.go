package model

type User struct {
	ID        uint `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique;not null;default:null"`
	Password  []byte

	Roles []Role `gorm:"many2many:user_roles"`
}

type DisplayUser struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

type CreateUserDTO struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}
