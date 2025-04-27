package domain

import (
	"time"
)

// Organization 组织领域模型
type Organization struct {
	ID          string    `json:"id"`          // 组织ID
	Code        string    `json:"code"`        // 组织代码
	ParentCode  string    `json:"parent_code"` // 父组织代码
	Name        string    `json:"name"`        // 组织名称
	UserID      string    `json:"user_id"`     // 用户ID
	Description string    `json:"description"` // 描述
	SortOrder   int       `json:"sort_order"`  // 排序序号
	CreatedAt   time.Time `json:"created_at"`  // 创建时间(ISO 8601格式)
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间(ISO 8601格式)
}

// DeleteOrganizationResponse represents the response for deleting an organization
type OrganizationSuccessResponse struct {
	Success bool `json:"success"`
}
type CreateOrganizationReq struct {
	ParentCode  string `json:"parent_code"` // 父组织代码
	Name        string `json:"name"`        // 组织名称
	Description string `json:"description"` // 描述
	SortOrder   int32  `json:"sort_order"`  // 排序序号
}

// GetUserOrganizationsResponse represents the response for getting user organizations
type ListOrganizationsResponse struct {
	Organizations []*Organization `json:"organizations"`
	Pagination
}

type AddItemRequest struct {
	OrganizationID string `json:"organization_id"`
	ItemID         string `json:"item_id"`
}
type RemoveItemRequest struct {
	OrganizationID string `json:"organization_id"`
	ItemID         string `json:"item_id"`
}

// AddMembersToOrganizationResponse represents the response for adding members to an organization
type AddItemsToOrganizationResponse struct {
	SuccessCount  int32    `json:"success_count"`
	FailureCount  int32    `json:"failure_count"`
	FailedItemIDs []string `json:"failed_item_ids"`
}

// RemoveMembersFromOrganizationResponse represents the response for removing members from an organization
type RemoveItemsFromOrganizationResponse struct {
	SuccessCount  int32    `json:"success_count"`
	FailureCount  int32    `json:"failure_count"`
	FailedItemIDs []string `json:"failed_item_ids"`
}
