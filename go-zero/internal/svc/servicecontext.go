// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package svc

import (
	"test-go/go-zero/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
