package tag

import (
	"time"
)

// Tag represents a tag in the system
type Tag struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	ItemCount   int32     `json:"item_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTagRequest represents the request to create a tag
type CreateTagRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// CreateTagResponse represents the response for creating a tag
type CreateTagResponse struct {
	Tag *Tag `json:"tag"`
}

// GetTagRequest represents the request to get a tag
type GetTagRequest struct {
	TagID  string `json:"tag_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

// GetTagResponse represents the response for getting a tag
type GetTagResponse struct {
	Tag *Tag `json:"tag"`
}

// UpdateTagRequest represents the request to update a tag
type UpdateTagRequest struct {
	TagID       string `json:"tag_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// UpdateTagResponse represents the response for updating a tag
type UpdateTagResponse struct {
	Tag *Tag `json:"tag"`
}

// DeleteTagRequest represents the request to delete a tag
type DeleteTagRequest struct {
	TagID  string `json:"tag_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

// DeleteTagResponse represents the response for deleting a tag
type DeleteTagResponse struct {
	Success bool `json:"success"`
}

// GetUserTagsRequest represents the request to get user tags
type GetUserTagsRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	Page        int32  `json:"page"`
	PageSize    int32  `json:"page_size"`
	SearchQuery string `json:"search_query"`
}

// GetUserTagsResponse represents the response for getting user tags
type GetUserTagsResponse struct {
	Tags       []*Tag `json:"tags"`
	Pagination struct {
		Page     int32 `json:"page"`
		PageSize int32 `json:"page_size"`
		Total    int32 `json:"total"`
		Pages    int32 `json:"pages"`
	} `json:"pagination"`
}

// AddTagsToItemsRequest represents the request to add tags to items
type AddTagsToItemsRequest struct {
	UserID  string   `json:"user_id" binding:"required"`
	Tags    []string `json:"tags" binding:"required"`
	ItemIDs []string `json:"item_ids" binding:"required"`
}

// AddTagsToItemsResponse represents the response for adding tags to items
type AddTagsToItemsResponse struct {
	SuccessCount   int32    `json:"success_count"`
	FailureCount   int32    `json:"failure_count"`
	FailedItemIDs  []string `json:"failed_item_ids"`
}

// RemoveTagsFromItemsRequest represents the request to remove tags from items
type RemoveTagsFromItemsRequest struct {
	UserID  string   `json:"user_id" binding:"required"`
	Tags    []string `json:"tags" binding:"required"`
	ItemIDs []string `json:"item_ids" binding:"required"`
}

// RemoveTagsFromItemsResponse represents the response for removing tags from items
type RemoveTagsFromItemsResponse struct {
	SuccessCount   int32    `json:"success_count"`
	FailureCount   int32    `json:"failure_count"`
	FailedItemIDs  []string `json:"failed_item_ids"`
}

// GetItemTagsRequest represents the request to get item tags
type GetItemTagsRequest struct {
	ItemID string `json:"item_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

// GetItemTagsResponse represents the response for getting item tags
type GetItemTagsResponse struct {
	ItemID string `json:"item_id"`
	Tags   []*Tag `json:"tags"`
} 