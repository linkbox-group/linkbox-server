package model

import "time"

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
