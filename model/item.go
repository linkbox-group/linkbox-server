package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/linkbox-group/linkbox-server/model/array"
	itemmodel "github.com/linkbox-group/linkbox-server/rpc-gen/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

// 项目模型
type Item struct {
	BaseModel
	UserID           string `gorm:"type:varchar(36);not null;index:idx_user_id;comment:用户ID" json:"userId"`
	ItemType         itemmodel.ItemType
	Title            string `gorm:"type:varchar(500);comment:标题" json:"title,omitempty"`
	Note             string `gorm:"type:text;comment:内容/文本" json:"note,omitempty"`
	URL              string `gorm:"type:varchar(2000);comment:链接地址" json:"url,omitempty"`
	ThumbnailURL     string `gorm:"type:varchar(2000);comment:缩略图地址" json:"thumbnailUrl,omitempty"`
	TagNames         array.Array
	OrganizationPath string

	CreatedAt CustomTime `json:"created_at"`
	UpdatedAt CustomTime `json:"updated_at"`
	DeletedAt gorm.DeletedAt
	// 关联
	User           User         `gorm:"foreignKey:UserID" json:"-"`
	OrganizationID string       `gorm:"type:varchar(36);not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organizations,omitempty"`
	Tags           []Tag        `gorm:"many2many:item_tag;" json:"tags,omitempty"`
}

func (Item) TableName() string {
	return "item"
}
func (i *Item) ConvertTo() *itemmodel.Item {
	return &itemmodel.Item{
		Id:               i.ID,
		UserId:           i.UserID,
		Type:             i.ItemType,
		Title:            i.Title,
		Url:              i.URL,
		Note:             i.Note,
		ThumbnailUrl:     i.ThumbnailURL,
		TagNames:         i.TagNames,
		OrganizationPath: i.OrganizationPath,
		OrganizationId:   i.OrganizationID,
		CreatedAt:        timestamppb.New(i.CreatedAt.Time()),
		UpdatedAt:        timestamppb.New(i.CreatedAt.Time()),
	}

}

type CustomTime time.Time

func (ct CustomTime) Value() (driver.Value, error) {
	t := time.Time(ct)
	// 不要返回 nil，而是返回实际的零值时间
	// 即使是零值时间也要返回有效的时间
	return t, nil
}

// Scan 实现 sql.Scanner 接口
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*ct = CustomTime(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*ct = CustomTime(v)
		return nil
	case []byte:
		t, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		*ct = CustomTime(t)
		return nil
	case string:
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		*ct = CustomTime(t)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into CustomTime", value)
	}
}
func (CustomTime) GormDataType() string {
	return "datetime"
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// 移除引号
	s := string(b)
	s = s[1 : len(s)-1]

	if s == "" {
		return nil
	}
	// 先尝试时间格式带空格的情况：2025-04-28 05:46:48
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		*ct = CustomTime(t)
		return nil
	}

	// 退化方案：尝试标准RFC3339格式
	t, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("unable to parse time: %v", err)
	}

	*ct = CustomTime(t)
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return json.Marshal(t.Format("2006-01-02 15:04:05"))
}

// 转换为普通 time.Time
func (ct CustomTime) Time() time.Time {
	return time.Time(ct)
}

// String 实现 fmt.Stringer 接口
func (ct CustomTime) String() string {
	return time.Time(ct).Format("2006-01-02 15:04:05")
}
