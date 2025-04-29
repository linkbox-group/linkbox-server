package domain

import "time"

type MessageResp struct {
	ID         string    `json:"id"`
	UserId     string    `json:"user_id"`
	Content    string    `json:"content"`
	SendTime   time.Time `json:"send_time"`
	SenderType string    `json:"sender_type"`
}
type ListMessagesResp struct {
	Messages []*MessageResp `json:"messages"`
	Pagination
}
