package core

import (
	logger "github.com/linkbox-group/linkbox-server/common/logs"
	"github.com/linkbox-group/linkbox-server/item/pkg/log"
)

const ServiceName = "organization"

func LoadLog() {
	log.RegisterLogger(logger.NewLogger(ServiceName))
}
