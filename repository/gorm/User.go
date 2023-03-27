package repo_gorm

import (
	"strings"

	"github.com/Ephilipz/fiberAuth/model"
	"gorm.io/gorm"
)

type UserGormRepo struct {
	db *gorm.DB
}

func NewUserGormRepo(db *gorm.DB) *UserGormRepo {
	return &UserGormRepo{
		db: db,
	}
}

func (r *UserGormRepo) Get(id uint) (model.User, error) {
	user := model.User{}
	err := r.db.Find(&user, id).Error
	return user, err
}

func (r *UserGormRepo) GetAll() ([]model.User, error) {
	users := []model.User{}
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserGormRepo) GetByEmail(email string) model.User {
	user := model.User{}
	r.db.Where("lower(email) = ?", strings.ToLower(email)).First(&user)
	return user
}

func (r *UserGormRepo) Create(user model.User) (uint, error) {
	defaultRoles := []model.Role{}
	r.db.Where("is_default = ?", true).Find(&defaultRoles)
	err := r.db.Create(&user).Association("Roles").Append(&defaultRoles)
	return user.ID, err
}

func (r *UserGormRepo) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserGormRepo) Update(user model.User) error {
	return r.db.Save(&user).Error
}

func (r *UserGormRepo) GetRoles(id uint) ([]model.Role, error) {
	roles := []model.Role{}
	user := model.User{Model: gorm.Model{ID: id}}
	err := r.db.Model(&user).Association("Roles").Find(&roles)
	return roles, err
}

func (r *UserGormRepo) UpdateRoles(userId uint, roleIds []uint) error {
	roles := []model.Role{}
	user := model.User{Model: gorm.Model{ID: userId}}
	err := r.db.Where("id in ?", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	return r.db.Model(&user).Association("Roles").Replace(&roles)
}
