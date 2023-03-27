package service

import (
	"strings"

	"github.com/Ephilipz/fiberAuth/model"
)

type roleRepository interface {
	Get(uint) model.Role
	GetAll() ([]model.Role, error)
	GetByName(string) model.Role
	Delete(uint) error
	Update(model.Role) error
	Create(model.Role) error
	CreateMultiple([]model.Role) error
}

type Role struct {
	repo roleRepository
}

func NewroleService(repo roleRepository) Role {
	return Role{
		repo: repo,
	}
}

func (s *Role) Get(id uint) model.Role {
	return s.repo.Get(id)
}

func (s *Role) GetAll() ([]model.Role, error) {
	return s.repo.GetAll()
}

func (s *Role) GetByName(name string) model.Role {
	return s.repo.GetByName(strings.Trim(strings.ToLower(name), " "))
}

func (s *Role) Create(role model.Role) error {
	return s.repo.Create(role)
}

func (s *Role) CreateMultiple(role []model.Role) error {
	return s.repo.CreateMultiple(role)
}

func (s *Role) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *Role) Update(role model.Role) error {
	return s.repo.Update(role)
}
