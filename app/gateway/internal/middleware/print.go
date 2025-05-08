package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

// PrintReq JWT token验证中间件 - 打印请求体
func PrintReq() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求方法和URL
		method := ctx.Request.Method
		url := ctx.Request.URL.String()

		// 打印请求基本信息
		fmt.Printf("[Request] %s %s\n", method, url)

		// 打印请求头信息
		fmt.Println("[Headers]:")
		for k, v := range ctx.Request.Header {
			fmt.Printf("  %s: %s\n", k, v)
		}

		// 如果有请求体，则读取并打印
		if ctx.Request.Body != nil {
			// 读取请求体
			bodyBytes, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				fmt.Printf("[Error reading body]: %v\n", err)
			} else {
				// 重要：将已读取的body放回，否则后续处理器将无法读取
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// 打印请求体内容
				fmt.Println("[Request Body]:")

				// 尝试将请求体解析为JSON格式化输出
				var prettyJSON bytes.Buffer
				err = json.Indent(&prettyJSON, bodyBytes, "", "  ")
				if err != nil {
					// 如果不是JSON格式，则直接打印原始内容
					fmt.Println(string(bodyBytes))
				} else {
					fmt.Println(prettyJSON.String())
				}
			}
		} else {
			fmt.Println("[Request Body]: Empty")
		}

		fmt.Println("----------------------------------")

		// 继续执行后续处理器
		ctx.Next()
	}
}
