package core

import (
	"github.com/linkbox-group/linkbox-server/ai/pkg/log"
	logger "github.com/linkbox-group/linkbox-server/common/logs"
)

const ServiceName = "Ai"

func LoadLog() {
	log.RegisterLogger(logger.NewLogger(ServiceName))
}
