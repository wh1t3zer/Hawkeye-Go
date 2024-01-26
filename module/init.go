package module

import (
	"fmt"

	"github.com/e421083458/golang_common/lib"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"google.golang.org/grpc"
)

// Instance ...
type Instance struct {
	ServiceName     string
	RegistryAddress string
	Services        []*registry.Service
}

// NewInstance ...
func NewInstance(name, addr string) (*Instance, error) {
	reg := consul.NewRegistry(registry.Addrs(lib.GetStringConf(addr)))
	array, err := reg.GetService(name)
	if err != nil {
		return nil, err
	}
	return &Instance{ServiceName: name, RegistryAddress: addr, Services: array}, nil
}

// Update 更新服务
func (r *Instance) Update() error {
	reg := consul.NewRegistry(registry.Addrs(lib.GetStringConf(r.RegistryAddress)))
	Services, err := reg.GetService(r.ServiceName)
	if err != nil {
		fmt.Printf("Failed Update Service, Info: %v\n", err)
		return err
	}
	r.Services = Services
	return nil
}

// GetNodeByRamdom 获取节点
func (r *Instance) GetNodeByRamdom() (node *registry.Node, err error) {
	next := selector.Random(r.Services)
	node, err = next()
	fmt.Println(node)
	if err != nil {
		fmt.Printf("Failed Get Next Node, Info: %v\n", err)
		return
	}
	return
}

// GetNodeByRoundRobin 获取节点
func (r *Instance) GetNodeByRoundRobin() (node *registry.Node, err error) {
	next := selector.RoundRobin(r.Services)
	node, err = next()
	if err != nil {
		fmt.Printf("Failed Get Next Node, Info: %v\n", err)
		return
	}
	return
}

// NewDefaultConn ...
func NewDefaultConn(name, addr string) (conn *grpc.ClientConn, err error) {
	// 1.连接注册中心获取服务节点
	user, err := NewInstance(name, addr)
	fmt.Println(user)
	if err != nil {
		fmt.Printf("Failed conn Registry center, info: %v\n", err)
		return
	}
	node, err := user.GetNodeByRamdom()
	if err != nil {
		fmt.Printf("Failed get node, info: %v\n", err)
		return
	}

	// 2.连接服务节点调用方法
	conn, err = grpc.Dial(node.Address, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Failed Call Service, Info: %v\n", err)
		return
	}
	return conn, nil
}
