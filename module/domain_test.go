package module

import (
	"context"
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
	fmt.Printf("%#v\n", node) // &registry.Node{Id:"q24cav79jo", Address:"172.31.50.249:57000", Metadata:map[string]string{}}
	// 2.连接服务节点调用方法
	conn, err := grpc.Dial(node.Address, grpc.WithInsecure())
	if err != nil {
		t.Errorf("Failed connect service node, info: %v\n", err)
		return
	}
	defer conn.Close()

	call := pb.NewVulClient(conn)
	req := &pb.PocRequest{
		Exploit: false,
		Target:  "172.31.50.252:7001",
		AssetId: "e9d20878683",
		// PocPlugins: []string{"id_2.py"}, // Weblogic_171023_wls_CVE_2017_10271_RCE.py  // 可写可不写
	}
	resp, err := call.Verify(context.Background(), req)
	if err != nil {
		t.Errorf("Failed call method(Verify), info: %#v\n", err)
		return
	}

	// 3.输出结果
	// result := &pb.PocResponse{
	// 	VerifyInfo:   resp.VerifyInfo,
	// 	ExploitInfo:  resp.ExploitInfo,
	// 	WebshellInfo: resp.WebshellInfo,
	// 	TrojanInfo:   resp.TrojanInfo,
	// }
	// data, _ := json.Marshal(result)
	fmt.Printf("%v\n", resp)
}
