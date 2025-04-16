package service

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/tag/internal/acl"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.TagServiceItf), new(*TagService)), NewTagService)
var _ acl.TagServiceItf = &TagService{}

type TagService struct {
	repo acl.TagRepositoryItf
}

func NewTagService(repo acl.TagRepositoryItf) *TagService {
	return &TagService{repo: repo}
}
