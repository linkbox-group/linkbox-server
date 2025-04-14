package delivery

import (
	"context"
	"github.com/linkbox-group/linkbox-server/rpc-gen/user"
)

// Register implements the UserDelivery interface.
func (s *UserDelivery) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserDelivery interface.
func (s *UserDelivery) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	// TODO: Your code here...
	return
}

// OAuthLogin implements the UserDelivery interface.
func (s *UserDelivery) OAuthLogin(ctx context.Context, req *user.OAuthLoginRequest) (resp *user.LoginResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserProfile implements the UserDelivery interface.
func (s *UserDelivery) GetUserProfile(ctx context.Context, req *user.UserIdRequest) (resp *user.GetUserProfileResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdateUserProfile implements the UserDelivery interface.
func (s *UserDelivery) UpdateUserProfile(ctx context.Context, req *user.UpdateUserProfileRequest) (resp *user.UpdateUserProfileResponse, err error) {
	// TODO: Your code here...
	return
}

// ChangePassword implements the UserDelivery interface.
func (s *UserDelivery) ChangePassword(ctx context.Context, req *user.ChangePasswordRequest) (resp *user.ChangePasswordResponse, err error) {
	// TODO: Your code here...
	return
}

// ForgotPassword implements the UserDelivery interface.
func (s *UserDelivery) ForgotPassword(ctx context.Context, req *user.ForgotPasswordRequest) (resp *user.ForgotPasswordResponse, err error) {
	// TODO: Your code here...
	return
}

// ResetPassword implements the UserDelivery interface.
func (s *UserDelivery) ResetPassword(ctx context.Context, req *user.ResetPasswordRequest) (resp *user.ResetPasswordResponse, err error) {
	// TODO: Your code here...
	return
}

// DeleteAccount implements the UserDelivery interface.
func (s *UserDelivery) DeleteAccount(ctx context.Context, req *user.DeleteAccountRequest) (resp *user.DeleteAccountResponse, err error) {
	// TODO: Your code here...
	return
}

// ListUsers implements the UserDelivery interface.
func (s *UserDelivery) ListUsers(ctx context.Context, req *user.ListUsersRequest) (resp *user.ListUsersResponse, err error) {
	// TODO: Your code here...
	return
}

// Logout implements the UserDelivery interface.
func (s *UserDelivery) Logout(ctx context.Context, req *user.LogoutRequest) (resp *user.LogoutResponse, err error) {
	// TODO: Your code here...
	return
}

// RefreshToken implements the UserDelivery interface.
func (s *UserDelivery) RefreshToken(ctx context.Context, req *user.RefreshTokenRequest) (resp *user.RefreshTokenResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserSubscription implements the UserDelivery interface.
func (s *UserDelivery) GetUserSubscription(ctx context.Context, req *user.UserIdRequest) (resp *user.GetUserSubscriptionResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdateUserSubscription implements the UserDelivery interface.
func (s *UserDelivery) UpdateUserSubscription(ctx context.Context, req *user.UpdateUserSubscriptionRequest) (resp *user.UpdateUserSubscriptionResponse, err error) {
	// TODO: Your code here...
	return
}
