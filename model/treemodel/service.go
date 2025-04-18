// tree_service.go
package treemodel

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// TreeService 提供树形模型的操作服务
type TreeService struct {
	DB *gorm.DB
}

// NewTreeService 创建树形服务实例
func NewTreeService(db *gorm.DB) *TreeService {
	return &TreeService{DB: db}
}

// extractTreeModel 从实体对象中提取TreeModel字段
// 返回TreeModel指针和原始反射值中的TreeModel字段
func extractTreeModel(entity interface{}) (*TreeModel, reflect.Value, error) {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if !val.IsValid() {
		return nil, reflect.Value{}, errors.New("无效的实体对象")
	}

	// 如果不是结构体
	if val.Kind() != reflect.Struct {
		return nil, reflect.Value{}, errors.New("实体必须是结构体类型")
	}

	// 遍历结构体的字段，找到TreeModel类型的字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		// 如果找到了TreeModel字段
		if fieldType.Type == reflect.TypeOf(TreeModel{}) {
			// 获取字段值
			tm := field.Interface().(TreeModel)
			return &tm, field, nil
		} else if fieldType.Type == reflect.TypeOf(&TreeModel{}) {
			// 处理指针类型
			if !field.IsNil() {
				tm := field.Interface().(*TreeModel)
				return tm, field, nil
			}
		}
	}

	return nil, reflect.Value{}, errors.New("实体中没有找到TreeModel字段")
}

// UpdateTreeLeaf 更新节点的叶子状态
func (s *TreeService) UpdateTreeLeaf(entity interface{}) error {
	treeModel, _, err := extractTreeModel(entity)
	if err != nil {
		return err
	}

	// 如果是根节点，无需更新
	if treeModel.Code == ROOT_ID {
		return nil
	}

	// 获取表名
	entityTable := entity.(schema.Tabler)
	tableName := entityTable.TableName()

	// 检查是否有子节点
	var childCount int64
	err = s.DB.Table(tableName).Where("parent_code = ?", treeModel.Code).Count(&childCount).Error
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	// 更新叶子状态
	leafStatus := LEAF
	if childCount > 0 {
		leafStatus = NON_LEAF
	}

	// 执行更新
	return s.DB.Table(tableName).
		Where("code = ?", treeModel.Code).
		Update("tree_leaf", leafStatus).
		Error
}

// ProcessTreeInfo 处理树形节点信息，在Save前调用
func (s *TreeService) ProcessTreeInfo(entity interface{}) error {

	treeModel, treeModelField, err := extractTreeModel(entity)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	// 表名
	entityTable := entity.(schema.Tabler)
	tableName := entityTable.TableName()
	// 如果是根节点
	if treeModel.ParentCode == "" || treeModel.ParentCode == ROOT_ID {
		treeModel.TreeLevel = 0
		treeModel.ParentCode = ROOT_ID
		treeModel.ParentCodes = ""
		treeModel.TreeNames = treeModel.Name

		// 根节点默认为非叶子节点
		treeModel.TreeLeaf = NON_LEAF
	} else {
		// 查找父节点
		logrus.Infoln(treeModel)
		var parentData TreeModel
		err := s.DB.Table(tableName).Where("code = ?", treeModel.ParentCode).Find(&parentData).Error
		if err != nil {
			logrus.Errorln(err)
			return err
		}
		// 获取父节点相关数据
		parentTreeLevel := parentData.TreeLevel
		parentParentCodes := parentData.ParentCodes
		parentTreeNames := parentData.TreeNames
		parentName := parentData.Name

		// 更新父节点的叶子状态为非叶子
		if parentData.TreeLeaf == LEAF {
			s.DB.Table(tableName).
				Where("code =?", treeModel.ParentCode).
				Update("tree_leaf", NON_LEAF)
		}

		// 计算当前节点层级
		treeModel.TreeLevel = parentTreeLevel + 1

		// 设置父节点编码链
		if parentParentCodes == "" {
			treeModel.ParentCodes = treeModel.ParentCode
		} else {
			treeModel.ParentCodes = parentParentCodes + CODE_SEPARATOR + treeModel.ParentCode
		}

		// 设置全名称
		if parentTreeNames == "" {
			treeModel.TreeNames = parentName + NAME_SEPARATOR + treeModel.Name
		} else {
			treeModel.TreeNames = parentTreeNames + NAME_SEPARATOR + treeModel.Name
		}

		// 默认设置为叶子节点
		treeModel.TreeLeaf = LEAF
	}

	// 将处理后的TreeModel值更新回原实体
	if treeModelField.CanSet() {
		if treeModelField.Kind() == reflect.Struct {
			treeModelField.Set(reflect.ValueOf(*treeModel))
		} else {
			// 指针类型
			treeModelField.Set(reflect.ValueOf(treeModel))
		}
	}

	return nil
}

// FindChildren 获取指定节点的直接子节点
func (s *TreeService) FindChildren(db *gorm.DB, dest interface{}, parentCode string) error {
	return db.Where("parent_code = ?", parentCode).
		Find(dest).Error
}

// FindByParentCodesLike 根据父节点编码链模糊查询
func (s *TreeService) FindByParentCodesLike(db *gorm.DB, dest interface{}, code string) error {
	if code == "" || code == ROOT_ID {
		return db.Find(dest).Error
	}

	return db.Where("parent_codes LIKE ?", "%"+code+"%").
		Or("code = ?", code).
		Find(dest).Error
}

// DeleteWithChildren 删除节点及其所有子节点
func (s *TreeService) DeleteWithChildren(db *gorm.DB, entity interface{}) error {
	treeModel, _, err := extractTreeModel(entity)
	if err != nil {
		return err
	}

	if treeModel.Code == ROOT_ID {
		return errors.New("不能删除根节点")
	}

	// 获取表名
	entityTable := entity.(schema.Tabler)
	tableName := entityTable.TableName()

	// 查找父节点
	var parentCode string
	db.Table(tableName).
		Where("code = ?", treeModel.Code).
		Select("parent_code").
		Scan(&parentCode)

	// 开始事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 删除子节点
		if err := tx.Table(tableName).
			Where("parent_codes LIKE ?", "%"+treeModel.Code+"%").
			Delete(nil).Error; err != nil {
			return err
		}

		// 删除当前节点
		if err := tx.Table(tableName).
			Where("code = ?", treeModel.Code).
			Delete(nil).Error; err != nil {
			return err
		}

		// 更新父节点的叶子状态
		if parentCode != "" && parentCode != ROOT_ID {
			var siblingCount int64
			tx.Table(tableName).
				Where("parent_code = ?", parentCode).
				Count(&siblingCount)

			if siblingCount == 0 {
				tx.Table(tableName).
					Where("code = ?", parentCode).
					Update("tree_leaf", LEAF)
			}
		}

		return nil
	})
}

// MoveNode 将节点移动到新的父节点下
func (s *TreeService) MoveNode(db *gorm.DB, entity interface{}, newParentCode string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		treeModel, _, err := extractTreeModel(entity)
		if err != nil {
			logrus.Errorln("提取TreeModel失败:", err)
			return err
		}

		// 验证移动有效性
		if treeModel.Code == newParentCode {
			return errors.New("不能将节点移动到自身下")
		}

		// 获取表名
		entityTable := entity.(schema.Tabler)
		tableName := entityTable.TableName()

		// 检查是否会造成循环引用
		if newParentCode != "" && newParentCode != ROOT_ID {
			var count int64
			if err := tx.Table(tableName).
				Where("code = ? AND parent_codes LIKE ?", newParentCode, "%"+treeModel.Code+"%").
				Count(&count).Error; err != nil {
				logrus.Errorln("检查循环引用时发生错误:", err)
				return err
			}
			if count > 0 {
				return errors.New("不能将节点移动到其子节点下")
			}
		}

		// 获取当前节点信息
		var origNode TreeModel
		if err := tx.Model(entity).Where("code = ?", treeModel.Code).First(&origNode).Error; err != nil {
			logrus.Errorln("获取原始节点信息失败:", err)
			return err
		}

		oldParentCode := origNode.ParentCode
		oldTreeLevel := origNode.TreeLevel
		oldParentCodes := origNode.ParentCodes
		oldTreeNames := origNode.TreeNames

		// 准备新节点信息
		var newNodeInfo struct {
			ParentCode  string
			TreeLevel   int
			ParentCodes string
			TreeNames   string
		}

		// 统一处理两种情况的节点信息计算
		if newParentCode == "" || newParentCode == ROOT_ID {
			// 移动到根节点下
			newNodeInfo.ParentCode = ROOT_ID
			newNodeInfo.TreeLevel = 0
			newNodeInfo.ParentCodes = ""
			newNodeInfo.TreeNames = treeModel.Name
		} else {
			// 移动到指定父节点下
			var newParent TreeModel
			if err := tx.Model(entity).Where("code = ?", newParentCode).First(&newParent).Error; err != nil {
				logrus.Errorln("获取新父节点信息失败:", err)
				return err
			}

			// 更新新父节点为非叶子节点
			if err := tx.Table(tableName).
				Where("code = ?", newParentCode).
				Update("tree_leaf", NON_LEAF).Error; err != nil {
				logrus.Errorln("更新父节点叶子状态失败:", err)
				return err
			}

			newNodeInfo.ParentCode = newParentCode
			newNodeInfo.TreeLevel = newParent.TreeLevel + 1

			// 构建父代码路径
			if newParent.ParentCodes == "" {
				newNodeInfo.ParentCodes = newParentCode
			} else {
				newNodeInfo.ParentCodes = newParent.ParentCodes + CODE_SEPARATOR + newParentCode
			}

			// 构建树名称路径
			if newParent.TreeNames == "" {
				newNodeInfo.TreeNames = newParent.Name + NAME_SEPARATOR + treeModel.Name
			} else {
				newNodeInfo.TreeNames = newParent.TreeNames + NAME_SEPARATOR + treeModel.Name
			}
		}

		// 更新当前节点
		updates := map[string]interface{}{
			"parent_code":  newNodeInfo.ParentCode,
			"tree_level":   newNodeInfo.TreeLevel,
			"parent_codes": newNodeInfo.ParentCodes,
			"tree_names":   newNodeInfo.TreeNames,
		}
		if err := tx.Table(tableName).Where("code = ?", treeModel.Code).Updates(updates).Error; err != nil {
			logrus.Errorln("更新当前节点失败:", err)
			return err
		}

		// 计算层级差异用于子节点更新
		levelDiff := newNodeInfo.TreeLevel - oldTreeLevel

		// 查找并更新子节点 (使用结构体切片)
		var descendants []TreeModel
		if err := tx.Table(tableName).
			Where("parent_codes LIKE ?", "%"+treeModel.Code+"%").
			Find(&descendants).Error; err != nil {
			logrus.Errorln("查找子节点失败:", err)
			return err
		}

		// 计算路径替换前缀
		oldPrefix := oldParentCodes
		newPrefix := newNodeInfo.ParentCodes

		if oldPrefix != "" {
			oldPrefix += CODE_SEPARATOR
		}
		if newPrefix != "" {
			newPrefix += CODE_SEPARATOR
		}

		// 更新子节点
		for _, desc := range descendants {
			// 计算新值
			newDescTreeLevel := desc.TreeLevel + levelDiff
			newDescParentCodes := strings.Replace(
				desc.ParentCodes,
				oldPrefix+treeModel.Code,
				newPrefix+treeModel.Code,
				1,
			)
			newDescTreeNames := strings.Replace(
				desc.TreeNames,
				oldTreeNames,
				newNodeInfo.TreeNames,
				1,
			)

			// 更新子节点
			if err := tx.Table(tableName).
				Where("code = ?", desc.Code).
				Updates(map[string]interface{}{
					"tree_level":   newDescTreeLevel,
					"parent_codes": newDescParentCodes,
					"tree_names":   newDescTreeNames,
				}).Error; err != nil {
				logrus.Errorln("更新子节点失败:", err, "节点代码:", desc.Code)
				return err
			}
		}

		// 更新原父节点的叶子状态
		if oldParentCode != "" && oldParentCode != ROOT_ID {
			var siblingCount int64
			if err := tx.Table(tableName).
				Where("parent_code = ?", oldParentCode).
				Count(&siblingCount).Error; err != nil {
				logrus.Errorln("检查原父节点子节点数失败:", err)
				return err
			}

			if siblingCount == 0 {
				if err := tx.Table(tableName).
					Where("code = ?", oldParentCode).
					Update("tree_leaf", LEAF).Error; err != nil {
					logrus.Errorln("更新原父节点叶子状态失败:", err)
					return err
				}
			}
		}

		return nil
	})
}

// FixTreeData 修复树形结构数据
func (s *TreeService) FixTreeData(db *gorm.DB, modelType interface{}) error {
	// 获取表名
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(modelType)
	tableName := stmt.Table

	// 重置所有根节点
	if err := db.Table(tableName).
		Where("parent_code = '' OR parent_code = ?", ROOT_ID).
		Updates(map[string]interface{}{
			"tree_level":   0,
			"parent_codes": "",
			"tree_leaf":    NON_LEAF,
		}).Error; err != nil {
		return err
	}

	// 获取所有根节点
	var rootNodes []map[string]interface{}
	if err := db.Table(tableName).
		Where("parent_code = '' OR parent_code = ?", ROOT_ID).
		Find(&rootNodes).Error; err != nil {
		return err
	}

	// 遍历根节点，修复其子树
	for _, root := range rootNodes {
		rootCode := root["code"].(string)
		rootName := root["name"].(string)

		db.Table(tableName).
			Where("code = ?", rootCode).
			Updates(map[string]interface{}{
				"tree_names": rootName,
			})

		// 递归修复子树
		s.fixTreeDataRecursive(db, tableName, rootCode, "", rootName)
	}

	return nil
}

// fixTreeDataRecursive 递归修复树形数据
func (s *TreeService) fixTreeDataRecursive(
	db *gorm.DB,
	tableName, parentCode, parentCodes, treeNames string,
) error {
	// 获取当前层次所有节点
	var nodes []map[string]interface{}
	if err := db.Table(tableName).
		Where("parent_code = ?", parentCode).
		Find(&nodes).Error; err != nil {
		return err
	}

	if len(nodes) == 0 {
		// 如果没有子节点，将当前节点设为叶子节点
		return db.Table(tableName).
			Where("code = ?", parentCode).
			Update("tree_leaf", LEAF).
			Error
	}

	// 计算父节点编码链
	newParentCodes := parentCodes
	if parentCode != "" && parentCode != ROOT_ID {
		if newParentCodes == "" {
			newParentCodes = parentCode
		} else {
			newParentCodes = newParentCodes + CODE_SEPARATOR + parentCode
		}
	}

	// 更新所有子节点
	for _, node := range nodes {
		nodeCode := node["code"].(string)
		nodeName := node["name"].(string)

		// 计算当前节点的层级
		treeLevel := 0
		if parentCode != "" && parentCode != ROOT_ID {
			treeLevel = strings.Count(newParentCodes, CODE_SEPARATOR) + 1
		}

		// 计算全名称
		newTreeNames := nodeName
		if treeNames != "" {
			newTreeNames = treeNames + NAME_SEPARATOR + nodeName
		}

		// 更新当前节点
		db.Table(tableName).
			Where("code = ?", nodeCode).
			Updates(map[string]interface{}{
				"tree_level":   treeLevel,
				"parent_codes": newParentCodes,
				"tree_names":   newTreeNames,
				"tree_leaf":    NON_LEAF, // 暂时设为非叶子，后续会根据是否有子节点更新
			})

		// 递归处理子节点
		s.fixTreeDataRecursive(db, tableName, nodeCode, newParentCodes, newTreeNames)
	}

	return nil
}
