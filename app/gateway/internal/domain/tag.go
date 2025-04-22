package domain

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

// TagResponse represents a tag response for API
type TagResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	ItemCount   int32     `json:"item_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TagListResponse represents a list of tags response for API
type TagListResponse struct {
	Tags       []TagResponse `json:"tags"`
	Pagination Pagination    `json:"pagination"`
}

// TagSuccessResponse represents a success response for tag operations
type TagSuccessResponse struct {
	Success bool `json:"success"`
}

// TagRelatedResponse represents structure for related tags response
type TagRelatedResponse struct {
	Tag          TagResponse `json:"tag"`
	Correlation  float32     `json:"correlation"`
	CoOccurrence int32       `json:"co_occurrence"`
}

// GetRelatedTagsFullResponse represents full response for related tags
type GetRelatedTagsFullResponse struct {
	TagID       string               `json:"tag_id"`
	RelatedTags []TagRelatedResponse `json:"related_tags"`
}

// ItemTagsResponse represents response for item tags
type ItemTagsResponse struct {
	ItemID string        `json:"item_id"`
	Tags   []TagResponse `json:"tags"`
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
	UserID string `json:"user_id" binding:"required"`
	PageRequest
	SearchQuery string `json:"search_query"`
}

// GetUserTagsResponse represents the response for getting user tags
type GetUserTagsResponse struct {
	Tags       []*Tag     `json:"tags"`
	Pagination Pagination `json:"pagination"`
}

// AddTagsToItemsRequest represents the request to add tags to items
type AddTagsToItemsRequest struct {
	UserID  string   `json:"user_id" binding:"required"`
	Tags    []string `json:"tags" binding:"required"`
	ItemIDs []string `json:"item_ids" binding:"required"`
}

// AddTagsToItemsResponse represents the response for adding tags to items
type AddTagsToItemsResponse struct {
	SuccessCount  int32    `json:"success_count"`
	FailureCount  int32    `json:"failure_count"`
	FailedItemIDs []string `json:"failed_item_ids"`
}

// RemoveTagsFromItemsRequest represents the request to remove tags from items
type RemoveTagsFromItemsRequest struct {
	UserID  string   `json:"user_id" binding:"required"`
	Tags    []string `json:"tags" binding:"required"`
	ItemIDs []string `json:"item_ids" binding:"required"`
}

// RemoveTagsFromItemsResponse represents the response for removing tags from items
type RemoveTagsFromItemsResponse struct {
	SuccessCount  int32    `json:"success_count"`
	FailureCount  int32    `json:"failure_count"`
	FailedItemIDs []string `json:"failed_item_ids"`
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

// GetRelatedTagsRequest represents the request to get related tags
type GetRelatedTagsRequest struct {
	TagID  string `json:"tag_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
	Limit  int32  `json:"limit"`
}

// RelatedTag represents a related tag with correlation information
type RelatedTag struct {
	Tag          *Tag    `json:"tag"`
	Correlation  float32 `json:"correlation"`
	CoOccurrence int32   `json:"co_occurrence"`
}

// GetRelatedTagsResponse represents the response for getting related tags
type GetRelatedTagsResponse struct {
	TagID       string       `json:"tag_id"`
	RelatedTags []RelatedTag `json:"related_tags"`
}

// MergeTagsRequest represents the request to merge tags
type MergeTagsRequest struct {
	UserID           string   `json:"user_id" binding:"required"`
	SourceTagIDs     []string `json:"source_tag_ids" binding:"required"`
	TargetTagID      string   `json:"target_tag_id" binding:"required"`
	DeleteSourceTags bool     `json:"delete_source_tags"`
}

// MergeTagsResponse represents the response for merging tags
type MergeTagsResponse struct {
	AffectedItems int32 `json:"affected_items"`
	TargetTag     *Tag  `json:"target_tag"`
}
