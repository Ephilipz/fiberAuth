package service

import (
	"errors"
	"strings"

	"github.com/Ephilipz/fiberAuth/model"
	"gorm.io/gorm"
)

type userRepository interface {
	Get(uint) (model.User, error)
	GetAll() ([]model.User, error)
	GetByEmail(string) (model.User, error)
	Create(model.User) (uint, error)
	Delete(uint) error
	Update(model.User) error
	UpdateRoles(uint, []uint) error
	GetRoles(uint) ([]model.Role, error)
}

type User struct {
	repo userRepository
}

func NewUserService(repo userRepository) User {
	return User{
		repo: repo,
	}
}

func (s *User) Get(id uint) (model.DisplayUser, error) {
	return mapResult(s.repo.Get(id))
}

func (s *User) GetAll() ([]model.DisplayUser, error) {
	return mapResults(s.repo.GetAll())
}

func (s *User) GetByEmail(email string) (model.DisplayUser, error) {
	result, err := s.repo.GetByEmail(email)
	return mapResult(result, err)
}

// @precondition : the email, firstname, lastname should not be empty.
// @precondition : password should be empty or longer than 7 characters
func (s *User) Create(userDTO model.CreateUserDTO) (uint, error) {
	if len(userDTO.FirstName) == 0 || len(userDTO.LastName) == 0 || len(userDTO.Email) == 0 {
		return 0, errors.New("firstname, lastname and email are required")
	}

	if len(userDTO.Password) > 0 && len(userDTO.Password) < 8 {
		return 0, errors.New("password must be at least 8 characters long")
	}

	password, err := GetHashedPasswordOrEmpty(userDTO.Password)
	if err != nil {
		return 0, err
	}

	user := model.User{
		Email:     strings.TrimSpace(strings.ToLower(userDTO.Email)),
		FirstName: strings.TrimSpace(userDTO.FirstName),
		LastName:  strings.TrimSpace(userDTO.LastName),
		Password:  password,
	}

	return s.repo.Create(user)
}

// @precondition : email cannot be empty
// @precondition : first and lastname cannot be empty
func (s *User) Update(userDTO model.UpdateUserDTO) error {
	if len(userDTO.FirstName) == 0 || len(userDTO.LastName) == 0 || len(userDTO.Email) == 0 {
		return errors.New("firstname, lastname and email are required")
	}
	user := model.User{
		Model:     gorm.Model{ID: userDTO.ID},
		FirstName: userDTO.FirstName,
		LastName:  userDTO.LastName,
		Email:     userDTO.Email,
	}
	return s.repo.Update(user)
}

func (s *User) ValidateCredentials(email string, password string) (bool, uint) {
	user, err := s.repo.GetByEmail(email)
	if err != nil || user.ID == 0 {
		return false, 0
	}
	isValid := CheckPasswordHash(password, user.Password)

	if !isValid {
		return false, 0
	}

	return true, user.ID
}

func (s *User) UpdateRoles(userId uint, roleIds []uint) error {
	return s.repo.UpdateRoles(userId, roleIds)
}

func (s *User) GetRoles(userId uint) ([]model.Role, error) {
	return s.repo.GetRoles(userId)
}

func mapResult(userModel model.User, err error) (model.DisplayUser, error) {
	if err != nil {
		return model.DisplayUser{}, err
	}
	return mapToDisplay(userModel), nil
}

func mapResults(userModels []model.User, err error) ([]model.DisplayUser, error) {
	if err != nil {
		return []model.DisplayUser{}, err
	}

	displayUsers := []model.DisplayUser{}
	for _, user := range userModels {
		displayUsers = append(displayUsers, mapToDisplay(user))
	}
	return displayUsers, nil
}

func mapToDisplay(user model.User) model.DisplayUser {
	return model.DisplayUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}
