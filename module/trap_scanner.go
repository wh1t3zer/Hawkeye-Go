package module

import (
	"context"
	"fmt"

	"github.com/e421083458/golang_common/lib"
	"github.com/wh1t3zer/Hawkeye-Go/dao"
	pb "github.com/wh1t3zer/Hawkeye-Go/micro/proto/grpc"
	"google.golang.org/grpc"
)

// TrapScanner ...
type TrapScanner struct {
	Conn *grpc.ClientConn
}

// NewTrapScanner ...
func NewTrapScanner() (*TrapScanner, error) {
	conn, err := NewDefaultConn(lib.GetStringConf("micro.aum.name"), lib.GetStringConf("micro.consul.registry_address"))
	if err != nil {
		return nil, err
	}
	return &TrapScanner{Conn: conn}, nil
}

// Run ...
func (v TrapScanner) Run(assetid, portid int64, target string, plugin *dao.TrapPluginInfo) ([]*dao.TrapInfo, error) {
	result := []*dao.TrapInfo{}
	// 2.调用方法
	call := pb.NewVulClient(v.Conn)

	req := &pb.TrapRequest{TargetList: []string{target}, TrapId: fmt.Sprintf("%v", plugin.ID), PluginText: plugin.Content}
	resp, err := call.Trap(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("Failed call method(Trap), info: %#v", err)
	}

	for _, item := range resp.Array {
		result = append(result, &dao.TrapInfo{
			AssetID: assetid, PortID: portid, PluginID: plugin.ID, Verify: item.Verify,
			TrapID: plugin.TrapID, Name: plugin.Name, Protocol: plugin.Protocol,
			AppName: plugin.AppName, HoneyPot: plugin.Honeypot, Desc: plugin.Desc,
		})
	}

	return result, nil
}
