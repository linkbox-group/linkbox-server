package domain

import (
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
	"github.com/linkbox-group/linkbox-server/rpc-gen/model"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

// ContentType represents the type of content
type ContentType string

const (
	ContentTypeLink     ContentType = "LINK"
	ContentTypeBookmark ContentType = "NOTE"
)

// Content represents a content item in the system
type Item struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Type             string    `json:"type"`
	Title            string    `json:"title"`
	URL              string    `json:"url"`
	Description      string    `json:"description"`
	ThumbnailURL     string    `json:"thumbnail_url"`
	TagNames         []string  `json:"tag_names"`
	OrganizationPath string    `json:"organization_path"`
	Tags             []string  `json:"tags"`
	OrganizationID   string    `json:"organization_id"`
	Note             string    `json:"note"`
	DeletedAt        time.Time `json:"deleted_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
type SearchItemsReq struct {
	Pagination *pagination.PaginationRequest
	Query      string `json:"query"`
	ItemType   string `json:"item_type"`
}

func (i *Item) Convert(item *model.Item) {
	i.ID = item.Id
	i.UserID = item.UserId
	i.Type = item.Type.String()
	i.Title = item.Title
	i.URL = item.Url
	i.Description = item.Description
	i.ThumbnailURL = item.ThumbnailUrl
	i.Tags = item.Tags
	i.OrganizationPath = item.OrganizationPath
	i.TagNames = item.TagNames
	i.OrganizationID = item.OrganizationId
	i.Note = item.Note
	i.DeletedAt = item.DeletedAt.AsTime()
	i.CreatedAt = item.CreatedAt.AsTime()
	i.UpdatedAt = item.UpdatedAt.AsTime()

}

// CreateContentRequest represents the request to create a content
type CreateItemRequest struct {
	Type           any      `json:"type"`
	URL            string   `json:"url"`
	Title          string   `json:"title" binding:"required"`
	Description    string   `json:"description"`
	ThumbnailURL   string   `json:"thumbnail_url"`
	Tags           []string `json:"tags"`
	OrganizationID string   `json:"organization_id"`
	Note           string   `json:"note"`
}

func (r *CreateItemRequest) ConvertTo() (i *item.CreateItemRequest) {

	i = &item.CreateItemRequest{
		Url:            r.URL,
		Title:          r.Title,
		Description:    r.Description,
		ThumbnailUrl:   r.ThumbnailURL,
		Tags:           r.Tags,
		OrganizationId: r.OrganizationID,
		Note:           r.Note,
	}
	logrus.Debug("itemtype", r.Type)
	logrus.Debug("typetype", reflect.ValueOf(r.Type).Type())
	if t, ok := r.Type.(int); ok {
		i.Type = model.ItemType(t)
		logrus.Debug("item type:", i.Type)
		return i
	}
	if t, ok := r.Type.(float64); ok {
		i.Type = model.ItemType(t)
		logrus.Debug("item type:", i.Type)
		return i
	}
	if t, ok := r.Type.(string); ok {
		i.Type = model.ItemType(model.ItemType_value[t])
		logrus.Debug("item type:", i.Type)
		return i
	}

	return i
}

// CreateContentResponse represents the response for creating a content
type CreateItemResponse struct {
	Item *Item `json:"item"`
}

// GetContentRequest represents the request to get a content
type GetItemRequest struct {
	ItemID string `json:"item_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

// GetContentResponse represents the response for getting a content
type GetItemResponse struct {
	Item *Item `json:"item"`
}

// UpdateContentRequest represents the request to update a content
type UpdateItemRequest struct {
	ItemID          string   `json:"item_id" binding:"required"`
	UserID          string   `json:"user_id" binding:"required"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	ThumbnailURL    string   `json:"thumbnail_url"`
	Tags            []string `json:"tags"`
	OrganizationIDs []string `json:"organization_ids"`
}

// UpdateContentResponse represents the response for updating a content
type UpdateItemResponse struct {
	Item *Item `json:"item"`
}

// DeleteContentRequest represents the request to delete a content
type DeleteItemRequest struct {
	ItemID string `json:"item_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

// DeleteContentResponse represents the response for deleting a content
type ItemSuccessResponse struct {
	Success bool `json:"success"`
}

// GetContentsByTagsRequest represents the request to get contents by tags
type GetItemsByTagsRequest struct {
	UserID   string   `json:"user_id" binding:"required"`
	Tags     []string `json:"tags" binding:"required"`
	Page     int32    `json:"page"`
	PageSize int32    `json:"page_size"`
}

// GetContentsByTagsResponse represents the response for getting contents by tags
type GetItemsByTagsResponse struct {
	Items []*Item `json:"items"`
	Pagination
}

// ItemListResponse represents the response for a list of items
type ItemListResponse struct {
	Items []*Item `json:"items"`
	Pagination
}
