package handler

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

// User ...
type User struct {
	ServiceName     string
	RegistryAddress string
	Services        []*registry.Service
}

// InitUser ...
func InitUser(name, addr string) *User {
	return &User{
		ServiceName:     name,
		RegistryAddress: addr,
	}
}

// RegistryConn 连接注册中心
func (u *User) RegistryConn() error {
	reg := consul.NewRegistry(registry.Addrs(lib.GetStringConf(u.RegistryAddress)))
	Services, err := reg.GetService(u.ServiceName)
	if err != nil {
		return err
	}
	u.Services = Services
	return nil
}

// GetNodeByRamdom 获取节点
func (u *User) GetNodeByRamdom() (*registry.Node, error) {
	next := selector.Random(u.Services)
	node, err := next()
	if err != nil {
		return nil, err
	}
	return node, nil
}
