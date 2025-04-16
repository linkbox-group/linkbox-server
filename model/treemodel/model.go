// tree_model.go
package treemodel

// TreeModel 可嵌入其他模型的树形结构基础模型
type TreeModel struct {
	Code        string `gorm:"column:code;size:64;not null;comment:节点编码" json:"code"`
	ParentCode  string `gorm:"column:parent_code;size:64;default:'0';comment:节点上级编码" json:"parentCode"`
	ParentCodes string `gorm:"column:parent_codes;size:1000;default:'';comment:节点所有上级编码" json:"parentCodes"`
	TreeLeaf    string `gorm:"column:tree_leaf;type:char(1);size:1;default:'1';comment:是否叶子节点(0:否 1:是)" json:"treeLeaf"`
	TreeLevel   int    `gorm:"column:tree_level;type:decimal(4,0);default:0;comment:节点层次级别(从0开始)" json:"treeLevel"`
	TreeNames   string `gorm:"column:tree_names;size:1000;default:'';comment:节点全名称(用/分隔)" json:"treeNames"`
	Name        string `gorm:"column:name;size:100;not null;comment:节点名称" json:"name"`
}

// 常量定义
const (
	CODE_SEPARATOR = ","
	NAME_SEPARATOR = "/"

	LEAF     = "1" // 叶子节点
	NON_LEAF = "0" // 非叶子节点

	ROOT_ID = "0" // 根节点ID
)

// TreeModelInterface 定义需要实现的接口，方便类型转换和反射操作
type TreeModelInterface interface {
	GetCode() string
	SetCode(code string)

	GetParentCode() string
	SetParentCode(code string)

	GetParentCodes() string
	SetParentCodes(codes string)

	GetTreeLeaf() string
	SetTreeLeaf(leaf string)

	GetTreeLevel() int
	SetTreeLevel(level int)

	GetTreeNames() string
	SetTreeNames(names string)

	GetName() string
	SetName(name string)
}

// 实现TreeModelInterface接口
func (tm *TreeModel) GetCode() string     { return tm.Code }
func (tm *TreeModel) SetCode(code string) { tm.Code = code }

func (tm *TreeModel) GetParentCode() string     { return tm.ParentCode }
func (tm *TreeModel) SetParentCode(code string) { tm.ParentCode = code }

func (tm *TreeModel) GetParentCodes() string      { return tm.ParentCodes }
func (tm *TreeModel) SetParentCodes(codes string) { tm.ParentCodes = codes }

func (tm *TreeModel) GetTreeLeaf() string     { return tm.TreeLeaf }
func (tm *TreeModel) SetTreeLeaf(leaf string) { tm.TreeLeaf = leaf }

func (tm *TreeModel) GetTreeLevel() int      { return tm.TreeLevel }
func (tm *TreeModel) SetTreeLevel(level int) { tm.TreeLevel = level }

func (tm *TreeModel) GetTreeNames() string      { return tm.TreeNames }
func (tm *TreeModel) SetTreeNames(names string) { tm.TreeNames = names }

func (tm *TreeModel) GetName() string     { return tm.Name }
func (tm *TreeModel) SetName(name string) { tm.Name = name }
