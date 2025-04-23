package log

import (
	"github.com/xyq777/holog"
)

func NewLogger(serviceName string) *holog.Logger {
	return holog.NewLogger(serviceName, holog.WithOutputStyle(holog.JSON))

}
