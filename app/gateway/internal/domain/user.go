package domain

import (
	"time"


)

type UserData struct {
	Id          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	DisplayName string   `json:"display_name"`
	AvatarUrl   string   `json:"avatar_url"`
	Roles       []string `json:"roles"`
	CreatedAt   string   `json:"created_at"`
}

type UserRegisterResp struct {
	UserData
}
type UserLoginResp struct {
	UserData
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type UserGetInfoResp struct {
	Id          string                 `json:"id"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email"`
	AvatarUrl   string                 `json:"avatar_url"`
	DisplayName string                 `json:"display_name"`
	Bio         string                 `json:"bio"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type UserUpdateResp struct {
	Id          string                 `json:"id"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email"`
	AvatarUrl   string                 `json:"avatar_url"`
	DisplayName string                 `json:"display_name"`
	Bio         string                 `json:"bio"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt    time.Time`json:"updated_at"`
}
type UserChangePasswordResp struct {
	Success bool `json:"success"`
}
type UserLogoutResp struct {
	Success bool   `json:"success"`
}

type TokenRefreshResp struct {
	AccessToken string `json:"access_token"`
	Success     bool   `json:"success"`
}
