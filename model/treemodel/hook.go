// tree_hooks.go
package treemodel

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TreeModelPlugin GORM插件，用于自动处理TreeModel的钩子
type TreeModelPlugin struct {
	TreeService *TreeService
}

// NewTreeModelPlugin 创建TreeModel插件
func NewTreeModelPlugin(db *gorm.DB) *TreeModelPlugin {
	return &TreeModelPlugin{
		TreeService: NewTreeService(db),
	}
}

// Name 实现gorm.Plugin接口
func (p *TreeModelPlugin) Name() string {
	return "TreeModelPlugin"
}

// Initialize 初始化插件，注册钩子函数
func (p *TreeModelPlugin) Initialize(db *gorm.DB) error {
	// 注册创建前的钩子
	db.Callback().Create().Before("gorm:create").Register("tree_model:before_create", p.BeforeCreate)

	// 注册创建后的钩子
	db.Callback().Create().After("gorm:create").Register("tree_model:after_create", p.AfterCreate)

	// 注册更新前的钩子

	return nil
}

// BeforeCreate 创建前处理TreeModel字段
func (p *TreeModelPlugin) BeforeCreate(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	// 检查表是否包含树形模型字段
	if hasTreeModelFields(db) {

		err := p.TreeService.ProcessTreeInfo(db.Statement.Dest)
		if err != nil {
			logrus.Errorf("TreeModel plugin: error processing tree info: %v", err)
			return
		}
	}
}

// AfterCreate 创建后更新树形关系
func (p *TreeModelPlugin) AfterCreate(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	if hasTreeModelFields(db) {

		err := p.TreeService.UpdateTreeLeaf(db.Statement.Dest)
		if err != nil {
			logrus.Errorf("TreeModel plugin: error updating tree leaf: %v", err)
			db.AddError(err)
		}
	}
}

// hasTreeModelFields 检查表是否包含树形结构需要的字段
func hasTreeModelFields(db *gorm.DB) bool {
	schema := db.Statement.Schema
	if schema == nil {
		return false
	}

	// 只检查最基本的几个字段即可
	requiredFields := []string{"code", "parent_code"}
	for _, field := range requiredFields {
		if _, ok := schema.FieldsByDBName[field]; !ok {
			return false
		}
	}

	return true
}
