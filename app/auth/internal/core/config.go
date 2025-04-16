package core

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

const (
	DefaultConfigFileType = "toml"
	DefaultConfigFilePath = "configs/%s/conf.toml"
)

var ProviderSet = wire.NewSet(NewContext, NewRedis)

func NewContext() context.Context {
	return context.Background()
}

type options struct {
	fileTypeOption string
	filePathOption string
}

type Option interface {
	apply(*options)
}

// o: immutable default config
var o = options{
	fileTypeOption: DefaultConfigFileType,
	filePathOption: DefaultConfigFilePath,
}

type fileTypeOption string

func (f fileTypeOption) apply(opt *options) {
	opt.fileTypeOption = string(f)
}
func WithFileTypeOption(fileType string) Option {
	return fileTypeOption(fileType)
}

type filePathOption string

func (f filePathOption) apply(opt *options) {
	opt.filePathOption = string(f)
}
func WithFilePathOption(fileType string) Option {
	return filePathOption(fileType)
}

func LoadConfig(opts ...Option) (err error) {
	var mode string
	for _, opt := range opts {
		opt.apply(&o)
	}
	rootDir, err := getProjectRoot()
	if err != nil {
		return fmt.Errorf("获取项目根目录失败: %w", err)
	}

	// 加载 .env 文件
	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {
		log.Println("未找到.env,使用默认环境变量")
	}
	viper.AutomaticEnv()
	if mode = viper.GetString("MODE"); mode == "" {
		mode = "test"
	}
	confPath := fmt.Sprintf(o.filePathOption, mode)
	if _, err := os.Stat(confPath); err != nil {
		return fmt.Errorf(
			"load config: %w", err)
	}

	viper.SetConfigFile(confPath)
	viper.SetConfigType(o.fileTypeOption)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf(
			"read config: %w", err)
	}
	return nil
}

// getProjectRoot 获取项目根目录
func getProjectRoot() (string, error) {
	// 从当前目录开始向上查找，直到找到 go.mod 文件
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("未找到 go.mod 文件")
		}
		dir = parent
	}
}
