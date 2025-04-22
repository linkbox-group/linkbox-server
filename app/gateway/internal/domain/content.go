package domain

import (
	"time"
)

// ContentType represents the type of content
type ContentType string

const (
	ContentTypeLink ContentType = "LINK"
)

// Content represents a content item in the system
type Content struct {
	ID            string          `json:"id"`
	UserID        string          `json:"user_id"`
	Type          ContentType     `json:"type"`
	URL           string          `json:"url"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	ThumbnailURL  string          `json:"thumbnail_url"`
	Tags          []string        `json:"tags"`
	OrganizationIDs []string        `json:"organization_ids"`
	IsFavorite    bool            `json:"is_favorite"`
	IsArchived    bool            `json:"is_archived"`
	IsPrivate     bool            `json:"is_private"`
	Metadata      ContentMetadata `json:"metadata"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	Note          string          `json:"note"`
	ReadCount     int32           `json:"read_count"`
}

// ContentMetadata represents metadata about the content
type ContentMetadata struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Author       string   `json:"author"`
	SiteName     string   `json:"site_name"`
	FaviconURL   string   `json:"favicon_url"`
	ContentType  string   `json:"content_type"`
	Keywords     []string `json:"keywords"`
	Language     string   `json:"language"`
	IsArticle    bool     `json:"is_article"`
}

// CreateContentRequest represents the request to create a content
type CreateContentRequest struct {
	UserID        string          `json:"user_id" binding:"required"`
	Type          ContentType     `json:"type" binding:"required"`
	URL           string          `json:"url" binding:"required"`
	Title         string          `json:"title" binding:"required"`
	Description   string          `json:"description"`
	ThumbnailURL  string          `json:"thumbnail_url"`
	Metadata      ContentMetadata `json:"metadata"`
	Tags          []string        `json:"tags"`
	CollectionIDs []string        `json:"collection_ids"`
	IsFavorite    bool            `json:"is_favorite"`
	IsPrivate     bool            `json:"is_private"`
	Note          string          `json:"note"`
}

// CreateContentResponse represents the response for creating a content
type CreateContentResponse struct {
	Content *Content `json:"content"`
}

// GetContentRequest represents the request to get a content
type GetContentRequest struct {
	ContentID string `json:"content_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}

// GetContentResponse represents the response for getting a content
type GetContentResponse struct {
	Content *Content `json:"content"`
}

// UpdateContentRequest represents the request to update a content
type UpdateContentRequest struct {
	ContentID     string          `json:"content_id" binding:"required"`
	UserID        string          `json:"user_id" binding:"required"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	ThumbnailURL  string          `json:"thumbnail_url"`
	Metadata      ContentMetadata `json:"metadata"`
	Tags          []string        `json:"tags"`
	CollectionIDs []string        `json:"collection_ids"`
	IsFavorite    bool            `json:"is_favorite"`
	IsArchived    bool            `json:"is_archived"`
	IsPrivate     bool            `json:"is_private"`
}

// UpdateContentResponse represents the response for updating a content
type UpdateContentResponse struct {
	Content *Content `json:"content"`
}

// DeleteContentRequest represents the request to delete a content
type DeleteContentRequest struct {
	ContentID string `json:"content_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}

// DeleteContentResponse represents the response for deleting a content
type DeleteContentResponse struct {
	Success bool `json:"success"`
}

// GetContentsByTagsRequest represents the request to get contents by tags
type GetContentsByTagsRequest struct {
	UserID   string   `json:"user_id" binding:"required"`
	Tags     []string `json:"tags" binding:"required"`
	Page     int32    `json:"page"`
	PageSize int32    `json:"page_size"`
}

// GetContentsByTagsResponse represents the response for getting contents by tags
type GetContentsByTagsResponse struct {
	Items      []*Content `json:"items"`
	Pagination struct {
		Total      int32 `json:"total"`
		Page       int32 `json:"page"`
		PageSize   int32 `json:"page_size"`
		TotalPages int32 `json:"total_pages"`
	} `json:"pagination"`
}
