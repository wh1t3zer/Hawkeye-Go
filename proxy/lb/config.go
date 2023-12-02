package lb

import (
	"fmt"

	"github.com/wh1t3zer/Hawkeye/proxy/zookeeper"
)

// Observer 观察员
type Observer interface {
	Update()
}

// LoadBalanceConf 配置主题
type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

// LoadBalanceZkConf zk配置
type LoadBalanceZkConf struct {
	observers    []Observer
	path         string
	zhHosts      []string
	confIPWeight map[string]string
	activeList   []string
	format       string
}

// Attach 新增观察员
func (s *LoadBalanceZkConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

// GetConf 获取配置
func (s *LoadBalanceZkConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIPWeight[ip]
		if !ok {
			weight = "50"
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

// WatchConf 监听配置, 通知观察员也更新
func (s *LoadBalanceZkConf) WatchConf() {
	zkManager := zookeeper.NewZkManager(s.zhHosts)
	zkManager.GetConnect()
	fmt.Println("watchConf")
	chanList, chanErr := zkManager.WatchServerListByPath(s.path)
	go func() {
		defer zkManager.Close()
		for {
			select {
			case changeErr := <-chanErr:
				fmt.Println("changeErr", changeErr)
			case changeList := <-chanList:
				fmt.Println("watch node changed")
				s.UpdateConf(changeList)
			}
		}
	}()
}

// UpdateConf 更新配置
func (s *LoadBalanceZkConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

// LoadBalanceConfObserver 负载均衡配置的观察员
type LoadBalanceConfObserver struct {
	ModuleConf *LoadBalanceZkConf
}
