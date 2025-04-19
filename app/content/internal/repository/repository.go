package repository

import (
	"context"
	"errors"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/content/internal/acl"
	"github.com/linkbox-group/linkbox-server/model"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.UserRepositoryItf), new(*Repository)), NewRepository)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) CreateItem(ctx context.Context, req *model.Item) (err error) {
	return r.db.Create(req).Error
}

func (r *Repository) GetItem(ctx context.Context, item *model.Item) (err error) {
	return r.db.
		Where("id = ? AND user_id = ?", item.ID, item.UserID).
		First(item).
		Error
}

func (r *Repository) UpdateItem(ctx context.Context, req *model.Item) (err error) {
	return r.db.
		Where("id = ? AND user_id = ?", req.ID, req.UserID).
		Updates(req).Error
}
func (r *Repository) DeleteItem(ctx context.Context, item *model.Item) (err error) {
	if item == nil {
		return errors.New("item is nil")
	}
	if item.ID == "" {
		return errors.New("item is empty")
	}
	if item.UserID == "" {
		return errors.New("item is empty")
	}
	result := r.db.Where("id = ? AND user_id = ?", item.ID, item.UserID).Delete(item)
	count := result.RowsAffected
	if result.Error != nil {
		return result.Error
	}
	if count == 0 {
		return errors.New("item not found")
	}
	return nil
	//return r.db.Where("id = ? AND user_id = ?", item.ID, item.UserID).Delete(item).Error
}
