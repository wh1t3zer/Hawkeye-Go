package lb

import (
	"errors"
	"fmt"
	"strings"
)

// RoundRobinBalance 随机负载
type RoundRobinBalance struct {
	curIndex int
	rss      []string
	//观察主体
	conf LoadBalanceConf
}

// Add 添加节点
func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

// Next 取下一个节点
func (r *RoundRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	// 直接+1 ,若+1后超出len, 则从头
	r.curIndex++
	if r.curIndex >= len(r.rss) {
		r.curIndex = 0
	}
	// r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

// Get 获取节点
func (r *RoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

// SetConf ...
func (r *RoundRobinBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

// Update 更新当前所有节点信息
func (r *RoundRobinBalance) Update() {
	if conf, ok := r.conf.(*LoadBalanceZkConf); ok {
		fmt.Println("update get conf:", conf.GetConf())
		r.rss = []string{}
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}
	if conf, ok := r.conf.(*LoadBalanceCheckConf); ok {
		fmt.Println("Update get conf:", conf.GetConf())
		r.rss = nil
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}
}
