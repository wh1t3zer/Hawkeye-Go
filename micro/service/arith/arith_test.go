package main

import (
	"context"
	"fmt"
	"testing"

	unit "github.com/wh1t3zer/Hawkeye-Go/micro/handler"
	pb "github.com/wh1t3zer/Hawkeye-Go/micro/proto/grpc"
	"google.golang.org/grpc"
)

func TestArith(t *testing.T) {
	// 1.连接注册中心获取服务节点
	user := unit.InitUser("test_python3", "127.0.0.1:8500")
	if err := user.RegistryConn(); err != nil {
		t.Errorf("Failed conn Registry center, info: %v\n", err)
		return
	}
	node, err := user.GetNodeByRamdom()
	if err != nil {
		t.Errorf("Failed get node, info: %v\n", err)
		return
	}
	fmt.Printf("current node info: %#v\n", node)
	// 2.连接服务节点调用方法
	conn, err := grpc.Dial(node.Address, grpc.WithInsecure())
	if err != nil {
		t.Errorf("Failed connect service node, info: %v\n", err)
		return
	}
	defer conn.Close()

	call := pb.NewArithClient(conn)
	resp, err := call.XiangJia(context.Background(), &pb.ArithRequest{Num1: 11, Num2: 22})
	if err != nil {
		t.Errorf("Failed call method(XiangJia), info: %#v\n", err)
		return
	}
	fmt.Printf("add_result=%v\n", resp.Result)
}
