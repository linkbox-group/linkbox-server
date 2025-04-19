package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/content/internal/acl"
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
	return r.db.
		Where("id = ? AND user_id = ?", req.ID, req.UserID).
		Updates(req).Error
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
func (r *Repository) GetItemsByTags(ctx context.Context, userID string, tagIDs []string, paginationReq *pagination.PaginationRequest) (items []model.Item, total int, err error) {
	// 计算分页参数
	limit, offset := int(paginationReq.GetPageSize()), int((paginationReq.Page-1)*paginationReq.PageSize)

	// 基础查询构建器，应用 context 和 user_id 过滤
	// 假设 model.Item 有 IsDeleted 字段用于软删除过滤
	baseQuery := r.db.WithContext(ctx).Model(&model.Item{}).Where("user_id = ? AND is_deleted = ?", userID, false)

	// 如果没有提供 tagIDs，则查询用户的所有未删除项目
	if len(tagIDs) == 0 {
		var totalCount int64
		// 计算总数
		if err = baseQuery.Count(&totalCount).Error; err != nil {
			return nil, 0, fmt.Errorf("counting user items failed: %w", err)
		}
		total = int(totalCount)

		// 获取分页数据
		if err = baseQuery.Order("created_at DESC").Offset(offset).Limit(limit).Preload("Tags").Find(&items).Error; err != nil {
			return nil, 0, fmt.Errorf("fetching user items failed: %w", err)
		}
		return items, total, nil
	}

	// --- 如果提供了 tagIDs ---

	// 构建带有 Joins, Group, Having 的查询，以查找包含所有指定标签的 items
	query := baseQuery.
		Joins("JOIN item_tags ON items.id = item_tags.item_id"). // 假设连接表名为 item_tags
		Where("item_tags.tag_id IN ?", tagIDs).
		Group("items.id").
		Having("COUNT(DISTINCT item_tags.tag_id) = ?", len(tagIDs))

	// 克隆查询用于计算总数（因为 Count() 会忽略 Group/Having，我们需要先获取 ID 列表）
	// 注意：直接在 Group/Having 查询上 Count 可能不准确，取决于数据库和 GORM 版本。
	// 更可靠的方法是先查询 ID，再 Count ID 数量。
	var ids []uint // 假设 Item 的 ID 是 uint 类型
	if err = query.Pluck("items.id", &ids).Error; err != nil {
		return nil, 0, fmt.Errorf("plucking item IDs by tags failed: %w", err)
	}
	total = len(ids) // 总数就是满足条件的 ID 数量

	// 如果没有找到匹配的 ID，直接返回空结果
	if total == 0 {
		return []model.Item{}, 0, nil
	}

	// 现在使用原始的 baseQuery（只过滤 user_id 和 is_deleted）
	// 并限制 ID 在我们找到的 ids 列表中，然后应用分页和 Preload
	if err = r.db.WithContext(ctx).Model(&model.Item{}).
		Where("id IN ?", ids).    // 使用上面找到的 ID 列表
		Order("created_at DESC"). // 保持排序一致性
		Offset(offset).           // 应用分页
		Limit(limit).             // 应用分页
		Preload("Tags").          // 预加载关联的 Tags
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("fetching items by tags failed: %w", err)
	}

	return items, total, nil
}

// ... 其他 repository 方法 ...
