package core

import (
	logger "github.com/linkbox-group/linkbox-server/common/logs"
	"github.com/linkbox-group/linkbox-server/user/pkg/log"
)

const ServiceName = "user"

func LoadLog() {
	log.RegisterLogger(logger.NewLogger(ServiceName))
}
