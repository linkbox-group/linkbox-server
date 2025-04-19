package service

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/user/internal/acl"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.UserServiceItf), new(*UserService)), NewConcreteUserUsecase)
