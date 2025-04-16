package repository

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/tag/internal/acl"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.TagRepositoryItf), new(*TagRepository)), NewTagRepository)

var _ acl.TagRepositoryItf = &TagRepository{}

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}
