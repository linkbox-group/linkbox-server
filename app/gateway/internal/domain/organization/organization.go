package organization

import (
	"time"
)

// Organization represents an organization in the system
type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	MemberCount int32     `json:"member_count"`
	ItemCount   int32     `json:"item_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateOrganizationRequest represents the request to create an organization
type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id" binding:"required"`
}

// CreateOrganizationResponse represents the response for creating an organization
type CreateOrganizationResponse struct {
	Organization *Organization `json:"organization"`
}

// GetOrganizationRequest represents the request to get an organization
type GetOrganizationRequest struct {
	ID      string `json:"id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

// GetOrganizationResponse represents the response for getting an organization
type GetOrganizationResponse struct {
	Organization *Organization `json:"organization"`
}

// UpdateOrganizationRequest represents the request to update an organization
type UpdateOrganizationRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"user_id" binding:"required"`
}

// UpdateOrganizationResponse represents the response for updating an organization
type UpdateOrganizationResponse struct {
	Organization *Organization `json:"organization"`
}

// DeleteOrganizationRequest represents the request to delete an organization
type DeleteOrganizationRequest struct {
	ID     string `json:"id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

// DeleteOrganizationResponse represents the response for deleting an organization
type DeleteOrganizationResponse struct {
	Success bool `json:"success"`
}

// GetUserOrganizationsRequest represents the request to get user organizations
type GetUserOrganizationsRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	Page        int32  `json:"page"`
	PageSize    int32  `json:"page_size"`
	SearchQuery string `json:"search_query"`
}

// GetUserOrganizationsResponse represents the response for getting user organizations
type GetUserOrganizationsResponse struct {
	Organizations []*Organization `json:"organizations"`
	Pagination struct {
		Page     int32 `json:"page"`
		PageSize int32 `json:"page_size"`
		Total    int32 `json:"total"`
		Pages    int32 `json:"pages"`
	} `json:"pagination"`
}

// AddMembersToOrganizationRequest represents the request to add members to an organization
type AddMembersToOrganizationRequest struct {
	OrganizationID string   `json:"organization_id" binding:"required"`
	UserID         string   `json:"user_id" binding:"required"`
	MemberIDs      []string `json:"member_ids" binding:"required"`
}

// AddMembersToOrganizationResponse represents the response for adding members to an organization
type AddMembersToOrganizationResponse struct {
	SuccessCount   int32    `json:"success_count"`
	FailureCount   int32    `json:"failure_count"`
	FailedMemberIDs []string `json:"failed_member_ids"`
}

// RemoveMembersFromOrganizationRequest represents the request to remove members from an organization
type RemoveMembersFromOrganizationRequest struct {
	OrganizationID string   `json:"organization_id" binding:"required"`
	UserID         string   `json:"user_id" binding:"required"`
	MemberIDs      []string `json:"member_ids" binding:"required"`
}

// RemoveMembersFromOrganizationResponse represents the response for removing members from an organization
type RemoveMembersFromOrganizationResponse struct {
	SuccessCount   int32    `json:"success_count"`
	FailureCount   int32    `json:"failure_count"`
	FailedMemberIDs []string `json:"failed_member_ids"`
}

// GetOrganizationMembersRequest represents the request to get organization members
type GetOrganizationMembersRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	UserID         string `json:"user_id" binding:"required"`
}

// GetOrganizationMembersResponse represents the response for getting organization members
type GetOrganizationMembersResponse struct {
	OrganizationID string   `json:"organization_id"`
	Members        []string `json:"members"`
} 