package repo_gorm

import (
	"strings"

	"github.com/Ephilipz/fiberAuth/model"
	"gorm.io/gorm"
)

type RoleGormRepo struct {
	db *gorm.DB
}

func NewRoleGormRepo(db *gorm.DB) *RoleGormRepo {
	return &RoleGormRepo{
		db: db,
	}
}

func (r *RoleGormRepo) Get(id uint) (model.Role, error) {
	role := model.Role{}
	err := r.db.Find(&role, id).Error
	return role, err
}

func (r *RoleGormRepo) GetAll() ([]model.Role, error) {
	roles := []model.Role{}
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *RoleGormRepo) GetByName(name string) model.Role {
	role := model.Role{}
	_ = r.db.Where("lower(name) = ?", strings.ToLower(name)).First(&role).Error
	return role
}

func (r *RoleGormRepo) Create(role model.Role) error {
	return r.db.Create(&role).Error
}

func (r *RoleGormRepo) CreateMultiple(roles []model.Role) error {
	return r.db.CreateInBatches(&roles, 100).Error
}

func (r *RoleGormRepo) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

func (r *RoleGormRepo) Update(role model.Role) error {
	return r.db.Save(&role).Error
}
