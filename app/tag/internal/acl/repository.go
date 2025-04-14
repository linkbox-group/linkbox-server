package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

type TagRepositoryItf interface {
	CreateTag(ctx context.Context, tag *model.Tag) (err error)
}
