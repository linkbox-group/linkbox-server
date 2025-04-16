package model

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
