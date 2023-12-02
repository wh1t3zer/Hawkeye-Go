package module

import (
	"context"
	"fmt"
	"log"

	"github.com/e421083458/golang_common/lib"
	"github.com/wh1t3zer/Hawkeye/dao"
	pb "github.com/wh1t3zer/Hawkeye/micro/proto/grpc"
	"google.golang.org/grpc"
)

// AssetScanner ...
type AssetScanner struct {
	IP       string
	PortList []string

	Conn *grpc.ClientConn
}

// NewAssetScanner ...
func NewAssetScanner(ip string, portlist []string) (*AssetScanner, error) {
	conn, err := NewDefaultConn(lib.GetStringConf("micro.aum.name"), lib.GetStringConf("micro.consul.registry_address"))
	if err != nil {
		return nil, err
	}
	return &AssetScanner{IP: ip, PortList: portlist, Conn: conn}, nil
}

// IPLocation ...
func (asr AssetScanner) IPLocation() (*dao.AssetInfo, error) {
	call := pb.NewHostClient(asr.Conn)
	req := &pb.LocRequest{Ip: asr.IP}
	resp, err := call.Location(context.Background(), req)
	if err != nil {
		log.Printf("Failed call method(IPLocation), info: %#v\n", err)
		return nil, err
	}
	return &dao.AssetInfo{AREA: resp.Area, ISP: resp.Isp, GPS: resp.Gps}, nil
}

// AliveList 网段存活列表
func (asr AssetScanner) AliveList(net string) ([]string, error) {
	call := pb.NewHostClient(asr.Conn)
	req := &pb.AlvRequest{Net: net}
	resp, err := call.Alive(context.Background(), req)
	if err != nil {
		log.Printf("Failed call method(AliveList), info: %#v\n", err)
		return nil, err
	}
	return resp.Hosts, nil
}

// HostDetail ...
func (asr AssetScanner) HostDetail() (*dao.AssetInfo, []*dao.PortInfo, error) {
	call := pb.NewHostClient(asr.Conn)
	req := &pb.DetlRequest{Ip: asr.IP, Ports: asr.PortList}
	resp, err := call.Detail(context.Background(), req)
	if err != nil {
		log.Printf("Failed call method(HostDetail), info: %#v\n", err)
		return nil, nil, err
	}
	var PortInfoList = []*dao.PortInfo{}
	for _, item := range resp.Array {
		fmt.Printf("%#v", item)
		PortInfoList = append(PortInfoList, &dao.PortInfo{
			Port: item.Port, State: item.State, Name: item.Name, Product: item.Product,
			Extrainfo: item.Extrainfo, Conf: item.Conf, Cpe: item.Cpe, Version: item.Version,
		})
	}

	return &dao.AssetInfo{OS: resp.Os, Vendor: resp.Vendor}, PortInfoList, nil
}
