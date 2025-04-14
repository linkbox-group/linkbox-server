package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

type TagServiceItf interface {
	CreateTagService(ctx context.Context, tag *model.Tag) (err error)
}
