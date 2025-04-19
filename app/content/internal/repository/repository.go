package repository

import (
	"context"
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
