package repository

import (
	"context"
	domain "github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/user/internal/acl"
	"gorm.io/gorm"
)

var _ acl.UserRepositoryItf = &MysqlUserRepo{}

type MysqlUserRepo struct {
	db *gorm.DB
}

func NewMysqlUserRepo(db *gorm.DB) *MysqlUserRepo {
	return &MysqlUserRepo{
		db: db,
	}
}

func (r *MysqlUserRepo) CreateUser(ctx context.Context, user *domain.User) (err error) {
	return r.db.WithContext(ctx).Model(&domain.User{}).Create(user).Error
}

func (r *MysqlUserRepo) FindUserByEmail(ctx context.Context, email string) (user *domain.User, err error) {
	user = &domain.User{}
	if err = r.db.
		WithContext(ctx).
		Model(&domain.User{}).
		Find(user, "email = ?", email).
		Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MysqlUserRepo) FindUserByID(ctx context.Context, id string) (user *domain.User, err error) {
	user = &domain.User{}
	if err = r.db.
		WithContext(ctx).
		Model(&domain.User{}).
		Find(user, "id = ?", id).
		Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MysqlUserRepo) UpdateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Model(user).Updates(user).Error
}

func (r *MysqlUserRepo) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Delete(&domain.User{}, id).Error
}
