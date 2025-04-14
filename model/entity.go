package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// 基础模型，包含通用字段
type BaseModel struct {
	ID        string    `gorm:"type:varchar(36);primaryKey;comment:UUID" json:"id"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// 为模型生成UUID
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == "" {
		base.ID = uuid.New().String()
	}
	return nil
}

// 用户模型
type User struct {
	BaseModel
	Email              string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_email;comment:邮箱" json:"email"`
	Phone              *string    `gorm:"type:varchar(20);uniqueIndex:idx_phone;comment:电话" json:"phone,omitempty"`
	PasswordHash       string     `gorm:"type:varchar(255);not null;comment:密码哈希" json:"-"`
	Username           string     `gorm:"type:varchar(50);not null;comment:用户名" json:"username"`
	AvatarURL          *string    `gorm:"type:varchar(255);comment:头像URL" json:"avatarUrl,omitempty"`
	Status             int8       `gorm:"type:tinyint;not null;default:1;comment:状态：1-正常，0-禁用" json:"status"`
	LastLoginAt        *time.Time `gorm:"type:datetime;comment:最后登录时间" json:"lastLoginAt,omitempty"`
	LoginCount         uint32     `gorm:"type:int unsigned;default:0;comment:登录次数" json:"loginCount"`
	RegistrationSource string     `gorm:"type:varchar(20);default:email;comment:注册来源: email,phone,wechat,github等" json:"registrationSource"`

	// 关联
	Subscriptions []UserSubscription `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`
	Collections   []Collection       `gorm:"foreignKey:UserID" json:"collections,omitempty"`
	Items         []Item             `gorm:"foreignKey:UserID" json:"items,omitempty"`
	Tags          []Tag              `gorm:"foreignKey:UserID" json:"tags,omitempty"`
	Events        []UserEvent        `gorm:"foreignKey:UserID" json:"events,omitempty"`
	PageVisits    []PageVisit        `gorm:"foreignKey:UserID" json:"pageVisits,omitempty"`
	FeatureUsages []FeatureUsage     `gorm:"foreignKey:UserID" json:"featureUsages,omitempty"`
}

func (User) TableName() string {
	return "user"
}

// 用户订阅模型
type UserSubscription struct {
	BaseModel
	UserID        string     `gorm:"type:varchar(36);not null;index:idx_user_id;comment:用户ID" json:"userId"`
	PlanType      string     `gorm:"type:varchar(20);not null;comment:订阅类型:free,pro,premium" json:"planType"`
	Status        string     `gorm:"type:varchar(20);not null;default:active;comment:状态:active,canceled,expired" json:"status"`
	StartDate     time.Time  `gorm:"type:datetime;not null;comment:开始时间" json:"startDate"`
	EndDate       *time.Time `gorm:"type:datetime;comment:结束时间" json:"endDate,omitempty"`
	PaymentMethod *string    `gorm:"type:varchar(50);comment:支付方式" json:"paymentMethod,omitempty"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (UserSubscription) TableName() string {
	return "user_subscription"
}

// 收藏夹模型
type Collection struct {
	BaseModel
	UserID        string     `gorm:"type:varchar(36);not null;index:idx_user_parent;comment:用户ID" json:"userId"`
	ParentID      *string    `gorm:"type:varchar(36);index:idx_user_parent;comment:父收藏夹ID" json:"parentId,omitempty"`
	Name          string     `gorm:"type:varchar(100);not null;comment:收藏夹名称" json:"name"`
	Description   *string    `gorm:"type:varchar(500);comment:描述" json:"description,omitempty"`
	IsDefault     bool       `gorm:"type:tinyint(1);not null;default:0;comment:是否默认收藏夹" json:"isDefault"`
	IsShared      bool       `gorm:"type:tinyint(1);not null;default:0;comment:是否共享" json:"isShared"`
	ShareCode     *string    `gorm:"type:varchar(32);index:idx_share_code;comment:分享码" json:"shareCode,omitempty"`
	ShareExpireAt *time.Time `gorm:"type:datetime;comment:分享过期时间" json:"shareExpireAt,omitempty"`
	SortOrder     int        `gorm:"type:int;not null;default:0;comment:排序顺序" json:"sortOrder"`
	ItemsCount    uint32     `gorm:"type:int unsigned;not null;default:0;comment:收藏项目数" json:"itemsCount"`

	// 关联
	User     User         `gorm:"foreignKey:UserID" json:"-"`
	Parent   *Collection  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Collection `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Items    []Item       `gorm:"many2many:collection_item;" json:"items,omitempty"`
}

func (Collection) TableName() string {
	return "collection"
}

// 收藏项目模型
type Item struct {
	BaseModel
	UserID          string     `gorm:"type:varchar(36);not null;index:idx_user_id;comment:用户ID" json:"userId"`
	ItemType        string     `gorm:"type:varchar(20);not null;index:idx_type;comment:类型:text,image,link,bookmark" json:"itemType"`
	Title           *string    `gorm:"type:varchar(500);comment:标题" json:"title,omitempty"`
	Content         *string    `gorm:"type:text;comment:内容/文本" json:"content,omitempty"`
	URL             *string    `gorm:"type:varchar(2000);comment:链接地址" json:"url,omitempty"`
	ImageURL        *string    `gorm:"type:varchar(2000);comment:图片地址" json:"imageUrl,omitempty"`
	ThumbnailURL    *string    `gorm:"type:varchar(2000);comment:缩略图地址" json:"thumbnailUrl,omitempty"`
	SourceDomain    *string    `gorm:"type:varchar(255);index:idx_domain;comment:来源网站域名" json:"sourceDomain,omitempty"`
	SourcePageTitle *string    `gorm:"type:varchar(500);comment:来源页面标题" json:"sourcePageTitle,omitempty"`
	IsFavorited     bool       `gorm:"type:tinyint(1);not null;default:0;comment:是否收藏" json:"isFavorited"`
	FavoritedAt     *time.Time `gorm:"type:datetime;comment:收藏时间" json:"favoritedAt,omitempty"`
	IsDeleted       bool       `gorm:"type:tinyint(1);not null;default:0;index:idx_user_deleted;comment:是否删除" json:"isDeleted"`
	DeletedAt       *time.Time `gorm:"type:datetime;index:idx_deleted_at;comment:删除时间" json:"deletedAt,omitempty"`

	// 关联
	User        User         `gorm:"foreignKey:UserID" json:"-"`
	Collections []Collection `gorm:"many2many:collection_item;" json:"collections,omitempty"`
	Tags        []Tag        `gorm:"many2many:item_tag;" json:"tags,omitempty"`
}

func (Item) TableName() string {
	return "item"
}

// 收藏夹项目关联模型
type CollectionItem struct {
	BaseModel
	CollectionID string `gorm:"type:varchar(36);not null;uniqueIndex:idx_collection_item;comment:收藏夹ID" json:"collectionId"`
	ItemID       string `gorm:"type:varchar(36);not null;uniqueIndex:idx_collection_item;index:idx_item_id;comment:项目ID" json:"itemId"`
	SortOrder    int    `gorm:"type:int;not null;default:0;comment:排序顺序" json:"sortOrder"`

	// 关联
	Collection Collection `gorm:"foreignKey:CollectionID" json:"-"`
	Item       Item       `gorm:"foreignKey:ItemID" json:"-"`
}

func (CollectionItem) TableName() string {
	return "collection_item"
}

// 标签模型
type Tag struct {
	BaseModel
	UserID   string  `gorm:"type:varchar(36);not null;uniqueIndex:idx_user_name;comment:用户ID" json:"userId"`
	Name     string  `gorm:"type:varchar(100);not null;uniqueIndex:idx_user_name;comment:标签名称" json:"name"`
	IsSystem bool    `gorm:"type:tinyint(1);not null;default:0;comment:是否系统标签" json:"isSystem"`
	Color    *string `gorm:"type:varchar(20);comment:标签颜色" json:"color,omitempty"`
	UseCount uint32  `gorm:"type:int unsigned;not null;default:0;comment:使用次数" json:"useCount"`

	// 关联
	User  User   `gorm:"foreignKey:UserID" json:"-"`
	Items []Item `gorm:"many2many:item_tag;" json:"items,omitempty"`
}

func (Tag) TableName() string {
	return "tag"
}

// 项目标签关联模型
type ItemTag struct {
	BaseModel
	ItemID string `gorm:"type:varchar(36);not null;uniqueIndex:idx_item_tag;comment:项目ID" json:"itemId"`
	TagID  string `gorm:"type:varchar(36);not null;uniqueIndex:idx_item_tag;index:idx_tag_id;comment:标签ID" json:"tagId"`

	// 关联
	Item Item `gorm:"foreignKey:ItemID" json:"-"`
	Tag  Tag  `gorm:"foreignKey:TagID" json:"-"`
}

func (ItemTag) TableName() string {
	return "item_tag"
}

// 用户事件模型
type UserEvent struct {
	BaseModel
	UserID     string  `gorm:"type:varchar(36);not null;index:idx_user_id;comment:用户ID" json:"userId"`
	EventType  string  `gorm:"type:varchar(50);not null;index:idx_event_type;comment:事件类型" json:"eventType"`
	EventName  string  `gorm:"type:varchar(100);not null;comment:事件名称" json:"eventName"`
	DeviceType *string `gorm:"type:varchar(20);comment:设备类型:web,mobile,app" json:"deviceType,omitempty"`
	Platform   *string `gorm:"type:varchar(20);comment:平台:ios,android,windows,macos,linux" json:"platform,omitempty"`
	Browser    *string `gorm:"type:varchar(50);comment:浏览器" json:"browser,omitempty"`
	IPAddress  *string `gorm:"type:varchar(50);comment:IP地址" json:"ipAddress,omitempty"`
	UserAgent  *string `gorm:"type:varchar(500);comment:用户代理" json:"userAgent,omitempty"`
	Properties *string `gorm:"type:json;comment:事件属性" json:"properties,omitempty"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (UserEvent) TableName() string {
	return "user_event"
}

// 页面访问模型
type PageVisit struct {
	BaseModel
	UserID      *string    `gorm:"type:varchar(36);index:idx_user_id;comment:用户ID,未登录为null" json:"userId,omitempty"`
	SessionID   string     `gorm:"type:varchar(100);not null;index:idx_session_id;comment:会话ID" json:"sessionId"`
	PageURL     string     `gorm:"type:varchar(500);not null;index:idx_page_url;comment:页面URL" json:"pageUrl"`
	PageTitle   *string    `gorm:"type:varchar(255);comment:页面标题" json:"pageTitle,omitempty"`
	Referrer    *string    `gorm:"type:varchar(500);comment:来源URL" json:"referrer,omitempty"`
	DeviceType  *string    `gorm:"type:varchar(20);comment:设备类型" json:"deviceType,omitempty"`
	Platform    *string    `gorm:"type:varchar(20);comment:平台" json:"platform,omitempty"`
	Browser     *string    `gorm:"type:varchar(50);comment:浏览器" json:"browser,omitempty"`
	IPAddress   *string    `gorm:"type:varchar(50);comment:IP地址" json:"ipAddress,omitempty"`
	EntryTime   time.Time  `gorm:"type:datetime;not null;index:idx_entry_time;comment:进入时间" json:"entryTime"`
	ExitTime    *time.Time `gorm:"type:datetime;comment:离开时间" json:"exitTime,omitempty"`
	Duration    *uint32    `gorm:"type:int unsigned;comment:停留时长(秒)" json:"duration,omitempty"`
	ScrollDepth *uint32    `gorm:"type:int unsigned;comment:滚动深度(%)" json:"scrollDepth,omitempty"`
	IsBounce    *bool      `gorm:"type:tinyint(1);comment:是否跳出" json:"isBounce,omitempty"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (PageVisit) TableName() string {
	return "page_visit"
}

// 功能使用统计模型
type FeatureUsage struct {
	BaseModel
	UserID      string     `gorm:"type:varchar(36);not null;uniqueIndex:idx_user_feature;comment:用户ID" json:"userId"`
	FeatureName string     `gorm:"type:varchar(100);not null;uniqueIndex:idx_user_feature;comment:功能名称" json:"featureName"`
	UsageCount  uint32     `gorm:"type:int unsigned;not null;default:0;comment:使用次数" json:"usageCount"`
	LastUsedAt  *time.Time `gorm:"type:datetime;comment:最后使用时间" json:"lastUsedAt,omitempty"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (FeatureUsage) TableName() string {
	return "feature_usage"
}
