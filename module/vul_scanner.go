package module

import (
	"context"
	"fmt"

	"github.com/e421083458/golang_common/lib"
	"github.com/wh1t3zer/Hawkeye-Go/dao"
	pb "github.com/wh1t3zer/Hawkeye-Go/micro/proto/grpc"
	"google.golang.org/grpc"
)

// VulScanner ...
type VulScanner struct {
	Conn *grpc.ClientConn
}

// NewVulScanner ...
func NewVulScanner() (*VulScanner, error) {
	conn, err := NewDefaultConn(lib.GetStringConf("micro.aum.name"), lib.GetStringConf("micro.consul.registry_address"))
	if err != nil {
		return nil, err
	}
	return &VulScanner{Conn: conn}, nil
}

// Pocsuite 漏洞验证、利用
func (v VulScanner) Pocsuite(assetID, portID int64, target string, pocArray []dao.PocPlugin) ([]*dao.VulInfo, error) {
	result := []*dao.VulInfo{}
	// 1.调用方法
	call := pb.NewVulClient(v.Conn)
	for _, plg := range pocArray {
		req := &pb.PocRequest{
			Exploit: true, VulId: fmt.Sprintf("%v", plg.ID), PocContent: plg.Content,
			AssetId: fmt.Sprintf("%v", portID), Target: target, //fmt.Sprintf("%v:%v", Host, v.Port),
		}
		resp, err := call.Verify(context.Background(), req)
		if err != nil {
			return nil, fmt.Errorf("Failed call method(Verify), info: %#v", err)
		}
		if resp.VerifyResult == "" {
			continue
		}
		result = append(result, &dao.VulInfo{
			AssetID:         assetID,
			PortID:          portID,
			PluginID:        plg.ID,
			VerifyURL:       resp.VerifyUrl,
			VerifyPayload:   resp.VerifyPayload,
			VerifyResult:    resp.VerifyResult,
			ExploitURL:      resp.ExploitUrl,
			ExploitPayload:  resp.ExploitPayload,
			ExploitResult:   resp.ExploitResult,
			WebshellURL:     resp.WebshellUrl,
			WebshellPayload: resp.WebshellPayload,
			WebshellResult:  resp.WebshellResult,
			TrojanURL:       resp.TrojanUrl,
			TrojanPayload:   resp.TrojanPayload,
			TrojanResult:    resp.TrojanResult,
		})
	}
	return result, nil
}
