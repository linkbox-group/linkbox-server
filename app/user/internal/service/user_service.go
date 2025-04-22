package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/auth"
	"github.com/linkbox-group/linkbox-server/rpc-gen/user"
	"github.com/linkbox-group/linkbox-server/user/internal/acl"
	"github.com/linkbox-group/linkbox-server/user/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/user/pkg/email"
	"github.com/linkbox-group/linkbox-server/user/pkg/encrypt"
	"github.com/linkbox-group/linkbox-server/user/pkg/regex"
	"github.com/sirupsen/logrus"

	"time"

	"gorm.io/gorm"
)

var (
	errUserAlreadyExist = errors.New("user already exist")
	errPasswordNotMatch = errors.New("password does not match")
)

const (
	TimeToExpire = 60 * time.Second
)

var _ acl.UserServiceItf = &UserService{}

type UserService struct {
	repo  acl.UserRepositoryItf
	cache *redis.Client
}

func NewConcreteUserUsecase(repo acl.UserRepositoryItf, cache *redis.Client) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

func (u *UserService) SendCode(ctx context.Context, sendTo string) (err error) {
	// 校验邮箱
	if regex.IsEmailInvalid(sendTo) {
		logrus.Errorln(err)
		return errors.New("invalid email")
	}

	// 是否 10 min 内重复发送验证码
	code := u.cache.Get(ctx, "email_code:"+sendTo).Val()
	if code != "" {
		logrus.Errorln(err)
		return errors.New("code already sent, please not resend within 10 min")
	}

	// 随机生成 6 位数字验证码
	code = email.RandomNumbers(6)

	// 将验证码存入 redis，10 min 后过期
	logrus.Info("用户 " + sendTo + " 验证码为 " + code)
	err = u.cache.Set(ctx, "email_code:"+sendTo, code, TimeToExpire).Err()
	if err != nil {
		logrus.Errorln(err)
		return errors.New("internal error: " + err.Error())
	}

	// 使用发送验证码
	go func() {
		err = email.SendVerificationCode(sendTo, code)
		if err != nil {
			logrus.Errorln(err)
		}
	}()

	return nil
}

func (u *UserService) RegisterUser(ctx context.Context, email, code, password string) (resp *user.RegisterResp, err error) {
	var (
		userModel *model.User
	)

	// 校验邮箱和密码格式
	err = regex.VerifyUser(email, password)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	// 检查用户邮箱是否已存在
	_, err = u.repo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Errorln(errUserAlreadyExist)
			return nil, errors.New("user already exist")
		}
		logrus.Errorln(err)
		return nil, err
	}

	// 校验验证码
	OriginCode := u.cache.Get(ctx, "email_code:"+email).Val()
	if code == "" || OriginCode != code {
		logrus.Errorln(err)
		return nil, errors.New("verification code is incorrect, please check again")
	}

	// 哈希密码
	passwordHash, err := encrypt.HashPassword(password)
	if err != nil {
		logrus.Errorln(err)
		return nil, errors.New("internal error: " + err.Error())
	}

	// 构建新用户
	userModel = &model.User{
		ID:           uuid.New().String(),
		Username:     "默认用户",
		Bio:          "该用户还没有没有填写个人简介",
		AvatarUrl:    "https://avatars.githubusercontent.com/u/204012462?s=48&v=4",
		Theme:        "light",
		RegisterDate: time.Now(),
		Email:        email,
		PasswordHash: passwordHash,
	}

	// 创建用户

	accessToken, err := rpc.AuthClient.GenerateAccessToken(ctx, &auth.GenerateTokenReq{
		Uid: userModel.ID,
	})
	if err != nil {
		logrus.Errorln(err)
		return nil, errors.New("internal error: " + err.Error())
	}
	refreshToken, err := rpc.AuthClient.GenerateRefreshToken(ctx, &auth.GenerateTokenReq{
		Uid: userModel.ID,
	})
	if err != nil {
		logrus.Errorln(err)
		return nil, errors.New("internal error: " + err.Error())
	}
	err = u.repo.CreateUser(ctx, userModel)
	if err != nil {
		logrus.Errorln(err)
		return nil, errors.New("internal error: " + err.Error())
	}

	// 返回用户信息
	return &user.RegisterResp{
		UserId:       userModel.ID,
		Username:     userModel.Username,
		Email:        userModel.Email,
		Bio:          userModel.Bio,
		Avatar:       userModel.AvatarUrl,
		Theme:        userModel.Theme,
		AccessToken:  accessToken.GetToken(),
		RefreshToken: refreshToken.GetToken(),
	}, nil
}

func (u *UserService) LoginUser(ctx context.Context, email, password string) (resp *user.LoginResp, err error) {
	// 校验邮箱和密码格式
	err = regex.VerifyUser(email, password)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	// 根据邮箱查找用户
	userModel, err := u.repo.FindUserByEmail(ctx, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Info(err)
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("internal error: %w", err)
	}

	// 校验密码
	if !encrypt.ComparePasswords(userModel.PasswordHash, password) {
		logrus.Info(errPasswordNotMatch)
		return nil, fmt.Errorf("password not match")
	}
	accessToken, err := rpc.AuthClient.GenerateAccessToken(ctx, &auth.GenerateTokenReq{
		Uid: userModel.ID,
	})
	if err != nil {
		logrus.Errorln(err)
		return nil, errors.New("internal error: " + err.Error())
	}
	refreshToken, err := rpc.AuthClient.GenerateRefreshToken(ctx, &auth.GenerateTokenReq{
		Uid: userModel.ID,
	})
	if err != nil {
		logrus.Errorln(err)
		return nil, errors.New("internal error: " + err.Error())
	}
	// 返回用户信息
	return &user.LoginResp{
		UserId:       userModel.ID,
		Username:     userModel.Username,
		Email:        userModel.Email,
		Bio:          userModel.Bio,
		Avatar:       userModel.AvatarUrl,
		Theme:        userModel.Theme,
		AccessToken:  accessToken.GetToken(),
		RefreshToken: refreshToken.GetToken(),
	}, nil
}

func (u *UserService) GetUser(ctx context.Context, id string) (resp *user.GetUserInfoResp, err error) {
	userModel, err := u.repo.FindUserByID(ctx, id)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	return &user.GetUserInfoResp{
		UserId:   userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Bio:      userModel.Bio,
		Avatar:   userModel.AvatarUrl,
		Theme:    userModel.Theme,
	}, nil
}

func (u *UserService) UpdatePassword(ctx context.Context, id string, oldPassword, newPassword string) (err error) {
	// 校验密码格式
	if regex.IsPasswordInvalid(newPassword) {
		logrus.Errorln(err)
		return errors.New("password invalid")
	}

	// 根据id查找用户
	userModel, err := u.repo.FindUserByID(ctx, id)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	// 校验密码
	if !encrypt.ComparePasswords(userModel.PasswordHash, oldPassword) {
		logrus.Info(errPasswordNotMatch)
		return fmt.Errorf("password not match")
	}

	// 更新密码
	password, err := encrypt.HashPassword(newPassword)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	userModel.PasswordHash = password
	if err := u.repo.UpdateUser(ctx, userModel); err != nil {
		logrus.Errorln(err)
		return err
	}

	return nil
}

// UpdateUserInfo 更新用户信息
func (u *UserService) UpdateUserInfo(ctx context.Context, req *user.UpdateUserInfoReq) (err error) {
	// 根据id查找用户
	userModel, err := u.repo.FindUserByID(ctx, req.UserId)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	if req.Username != "" {
		userModel.Username = req.Username
	}
	if req.Bio != "" {
		userModel.Bio = req.Bio
	}
	if req.Avatar != "" {
		userModel.AvatarUrl = req.Avatar
	}
	if req.Theme != "" {
		userModel.Theme = req.Theme
	}

	// 更新用户信息
	if err := u.repo.UpdateUser(ctx, userModel); err != nil {
		logrus.Errorln(err)
		return err
	}

	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, id string) (err error) {
	if err := u.repo.DeleteUser(ctx, id); err != nil {
		logrus.Errorln(err)
		return err
	}

	if err != nil {
		logrus.Errorln(err)
		return err
	}

	return nil
}
