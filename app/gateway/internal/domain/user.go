package domain

type UserRegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}
type UserRegisterResp struct {
	UserId       string `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Avatar       string `json:"avatar"`
	Bio          string `json:"bio"`
	Theme        string `json:"theme"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type UserLoginResp struct {
	UserId       string `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Avatar       string `json:"avatar"`
	Bio          string `json:"bio"`
	Theme        string `json:"theme"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type UserGetInfoResp struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Theme    string `json:"theme"`
}
type UserUpdateResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
type UserChangePasswordResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
type UserLogoutResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
