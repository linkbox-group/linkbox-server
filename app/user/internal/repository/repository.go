package repository

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/user/internal/acl"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.UserRepositoryItf), new(*MysqlUserRepo)), NewMysqlUserRepo)
