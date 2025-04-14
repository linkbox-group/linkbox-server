package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/linkbox-group/linkbox-server/model"
)

var (
	ErrDbCreateTagFailed = errors.New("create tag failed")
)

func (s *TagService) CreateTagService(ctx context.Context, tag *model.Tag) (err error) {
	err = s.repo.CreateTag(ctx, tag)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrDbCreateTagFailed, err)
	}
	return nil
}
