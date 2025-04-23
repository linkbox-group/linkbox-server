package log

import "github.com/xyq777/holog"

var logger *holog.Logger

func Log() *holog.Logger {
	return logger
}
func RegisterLogger(exLogger *holog.Logger) {
	logger = exLogger

}
