package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	unit "github.com/wh1t3zer/Hawkeye/micro/handler"
	pb "github.com/wh1t3zer/Hawkeye/micro/proto/grpc"
	"google.golang.org/grpc"
)

// TestVerify Poc 漏洞验证及漏洞扫描
func TestVerify(t *testing.T) {
	// 1.连接注册中心获取服务节点
	user := unit.InitUser("aquaman", "172.31.50.249:8500")
	if err := user.RegistryConn(); err != nil {
		t.Errorf("Failed conn Registry center, info: %v\n", err)
		return
	}
	node, err := user.GetNodeByRamdom()
	if err != nil {
		t.Errorf("Failed get node, info: %v\n", err)
		return
	}
	fmt.Printf("node: %#v\n", node) // &registry.Node{Id:"q24cav79jo", Address:"172.31.50.249:57000", Metadata:map[string]string{}}
	// 2.连接服务节点调用方法
	conn, err := grpc.Dial(node.Address, grpc.WithInsecure())
	if err != nil {
		t.Errorf("Failed connect service node, info: %v\n", err)
		return
	}
	defer conn.Close()

	// [+] 漏洞类
	call := pb.NewVulClient(conn)
	req := &pb.PocRequest{Exploit: true, Target: "172.31.50.252:8081", AssetId: "hxyx"} // target必须带端口
	resp, err := call.Verify(context.Background(), req)
	if err != nil {
		t.Errorf("Failed call method(Verify), info: %#v\n", err)
		return
	}

	// 3.输出结果
	data, _ := json.Marshal(resp)
	fmt.Printf("%s\n", data)
}
