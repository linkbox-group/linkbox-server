package mysql_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/item/internal/acl"
	"github.com/linkbox-group/linkbox-server/item/pkg/log"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.UserRepositoryItf), new(*Repository)), NewRepository)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) CreateItem(ctx context.Context, req *model.Item) (err error) {
	return r.db.Create(req).Error
}

func (r *Repository) GetItem(ctx context.Context, item *model.Item) (err error) {
	return r.db.
		Where("id = ? AND user_id = ?", item.ID, item.UserID).
		First(item).
		Error
}

func (r *Repository) UpdateItem(ctx context.Context, req *model.Item) (err error) {
	//logrus.Infoln(req)
	res := r.db.
		Where("id = ? AND user_id = ?", req.ID, req.UserID).
		Updates(req)
	if res.RowsAffected == 0 {

		log.Log().Error("update item failed")
		return gorm.ErrRecordNotFound
	}
	return res.Error
}
func (r *Repository) DeleteItem(ctx context.Context, item *model.Item) (err error) {
	if item == nil {
		return errors.New("item is nil")
	}
	if item.ID == "" {
		return errors.New("item is empty")
	}
	if item.UserID == "" {
		return errors.New("item is empty")
	}
	result := r.db.Where("id = ? AND user_id = ?", item.ID, item.UserID).Delete(item)
	count := result.RowsAffected
	if result.Error != nil {
		return result.Error
	}
	if count == 0 {
		return errors.New("item not found")
	}
	return nil
	//return r.db.Where("id = ? AND user_id = ?", item.ID, item.UserID).Delete(item).Error
}
func (r *Repository) GetItemsByTags(ctx context.Context, userID string, tagNames []string, paginationReq *pagination.PaginationRequest) (items []model.Item, total int, err error) {
	// 计算分页参数
	limit, offset := int(paginationReq.GetPageSize()), int((paginationReq.Page-1)*paginationReq.PageSize)

	// 初始化查询构建器，应用 context 和 user_id 过滤
	// 注意修正：deleted_at IS NULL 表示未删除的项目（基于软删除机制）
	db := r.db.WithContext(ctx).Model(&model.Item{}).Where("user_id = ? AND deleted_at IS NULL", userID)

	// 记录总数查询器
	countQuery := db.Session(&gorm.Session{})

	// 如果提供了标签名称，使用子查询构建高效的标签筛选
	if len(tagNames) > 0 {
		// 获取标签 ID 不必作为单独的查询
		// 构建参数化子查询，找到包含所有指定标签的项目 ID
		subQuery := r.db.WithContext(ctx).Table("item_tag").
			Select("item_id").
			Joins("JOIN tag ON item_tag.tag_id = tag.id").
			Where("tag.user_id = ? AND tag.name IN ?", userID, tagNames).
			Group("item_id").
			Having("COUNT(DISTINCT tag.name) = ?", len(tagNames))

		// 将子查询应用到主查询和计数查询
		db = db.Where("id IN (?)", subQuery)
		countQuery = countQuery.Where("id IN (?)", subQuery)
	}

	// 执行计数查询
	var totalCount int64
	if err = countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("counting items failed: %w", err)
	}
	total = int(totalCount)

	// 如果没有结果，提前返回空数组
	if total == 0 {
		return []model.Item{}, 0, nil
	}

	// 获取分页后的数据，并预加载关联的标签
	if err = db.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Preload("Tags").
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("fetching items failed: %w", err)
	}

	return items, total, nil
}

func (r *Repository) GetItemsByOrganization(ctx context.Context, userID string, organizationID string, pageNum int, pageSize int) (items []model.Item, total int, err error) {

	// 参数验证
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10 // 默认页大小
	}

	// 计算分页参数
	limit, offset := pageSize, (pageNum-1)*pageSize

	// 创建查询构建器，使用 Model 而不是 Table 以获得更好的类型安全性
	query := r.db.WithContext(ctx).Model(&model.Item{}).Where("user_id = ? AND deleted_at IS NULL", userID) // 修正了软删除条件

	// 添加组织ID过滤
	if organizationID != "" {
		query = query.Where("organization_id = ?", organizationID)
	}

	// 计算总数
	var totalCount int64
	countQuery := query.Session(&gorm.Session{}) // 创建查询的副本用于计数
	if err = countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("counting organization items failed: %w", err)
	}
	total = int(totalCount)

	// 如果没有结果，提前返回
	if total == 0 {
		return []model.Item{}, 0, nil
	}

	// 获取分页数据
	if err = query.
		Select("item.*"). // 明确指定只选择 item 表的字段，避免列名冲突
		Order("item.created_at DESC").
		Offset(offset).
		Limit(limit).
		Preload("Tags", "deleted_at IS NULL"). // 只加载未删除的标签
		Find(&items).Error; err != nil {
		log.Log().Error(err.Error())
		return nil, 0, fmt.Errorf("fetching organization items failed: %w", err)
	}

	return items, total, nil
}

func (r *Repository) SearchItemsByTitle(ctx context.Context, userID string, query string, pageNum int, pageSize int) (items []model.Item, total int, err error) {
	// 计算分页参数
	limit, offset := pageSize, (pageNum-1)*pageSize
	// 基础查询构建器，使用连接表查询
	baseQuery := r.db.WithContext(ctx).
		Table("item").
		Where("item.user_id =? AND item.deleted_at is not null", userID)
	// 添加组织ID过滤
	if query != "" {
		baseQuery = baseQuery.Where("item.title LIKE ?", "%"+query+"%")
	}

	// 计算总数
	var totalCount int64
	if err = baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("counting user item failed: %w", err)
	}
	total = int(totalCount)
	// 获取分页数据
	if err = baseQuery.
		Order("item.created_at DESC").
		Offset(offset).
		Limit(limit).
		Preload("Tags").
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("fetching user items failed: %w", err)
	}
	return items, total, nil

}

func (r *Repository) RecoverItemsBatch(ctx context.Context, userID string, ids []string) error {
	if len(ids) == 0 {
		println("No IDs provided for recovery.")
		return nil
	}
	result := r.db.WithContext(ctx).Unscoped().
		Model(&model.Item{}).
		Where("user_id = ? AND id IN ? AND deleted_at IS NOT NULL", userID, ids).
		Update("deleted_at", nil)
	if result.Error != nil {
		return fmt.Errorf("failed to recover items: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no items found for recovery")
	}
	return nil
}
func (r *Repository) DeleteItemsBatch(ctx context.Context, userID string, ids []string) error {
	if len(ids) == 0 {
		println("No IDs provided for deletion.")
		return nil
	}
	result := r.db.WithContext(ctx).
		Where("user_id =? AND id IN?", userID, ids).
		Delete(&model.Item{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete items: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no items found for deletion")
	}
	return nil
}
