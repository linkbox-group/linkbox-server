package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/common/ecode"
	"github.com/linkbox-group/linkbox-server/gateway/internal/domain"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/rpc-gen/auth"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			domain.Error(ctx, ecode.ErrTokenInvalid, "Token 为空")
			return
		}

		// 检查 Authorization 头是否以 "Bearer " 开头
		const bearerPrefix = "Bearer "
		if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
			// 去掉 "Bearer " 前缀，获取实际的 Token
			tokenString := authHeader[len(bearerPrefix):]

			// 调用auth服务验证token
			tokenInfo, err := rpc.AuthClient.VerifyAccessToken(context.Background(), &auth.TokenRequest{
				Token: tokenString,
			})
			if err != nil {
				domain.Error(ctx, ecode.ErrTokenInvalid, "Token 验证错误")
				return
			}

			// 将用户ID存储在上下文中
			ctx.Set("userId", tokenInfo.Uid)
			ctx.Next()
		} else {
			domain.Error(ctx, ecode.ErrTokenInvalid, "Token 格式错误")
			return
		}
	}
}
