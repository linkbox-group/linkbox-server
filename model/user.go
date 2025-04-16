package model

import (
	"time"
)

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
	Organizations []Organization     `gorm:"foreignKey:UserID" json:"collections,omitempty"`
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
