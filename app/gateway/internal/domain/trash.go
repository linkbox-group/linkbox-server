package domain

import (
	"github.com/linkbox-group/linkbox-server/rpc-gen/model"
	"time"
)

type TrashItem struct {
	DeletedAt time.Time `json:"deleted_at"`
	ExpiredAt time.Time `json:"expired_at"`
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UserID    string    `json:"user_id"`
}
type ListTrashItemResp struct {
	Items      []*TrashItem `json:"items"`
	Page       int32        `json:"page"`
	PageSize   int32        `json:"page_size"`
	Total      int32        `json:"total"`
	TotalPages int32        `json:"total_pages"`
}

func (i *TrashItem) Convert(ti *model.TrashItem) {
	i.ID = ti.Id
	i.Title = ti.Title
	i.UserID = ti.UserId
	i.ExpiredAt = ti.ExpiredAt.AsTime()
	i.DeletedAt = ti.DeletedAt.AsTime()
}

type RecoveryTrashItemReq struct {
	ItemID string `json:"item_id"`
}

type RecoveryTrashItemResp struct {
	Success bool `json:"success"`
}
