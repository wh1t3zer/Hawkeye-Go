package module

import (
	"context"
	"fmt"
	"strings"

	"github.com/e421083458/golang_common/lib"
	"github.com/wh1t3zer/Hawkeye-Go/dao"
	pb "github.com/wh1t3zer/Hawkeye-Go/micro/proto/grpc"
	"google.golang.org/grpc"
)

// WebScanner ...
type WebScanner struct {
	Conn *grpc.ClientConn
}

// NewWebScanner ...
func NewWebScanner() (*WebScanner, error) {
	conn, err := NewDefaultConn(lib.GetStringConf("micro.aum.name"), lib.GetStringConf("micro.consul.registry_address"))
	if err != nil {
		return nil, err
	}
	return &WebScanner{Conn: conn}, nil
}

// WebSpider 爬虫
func (spi WebScanner) WebSpider(portID int64, Host string, Port int32) (*dao.WebInfo, error) {
	// 2.调用方法
	call := pb.NewWebScrapClient(spi.Conn)
	req := &pb.SpiRequest{Host: Host, Port: Port}
	resp, err := call.Spider(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("Failed call method(Resolv), info: %v", err)
	}
	if resp.StartUrl == "" && resp.Title == "" && resp.Server == "" && len(resp.RouteList) < 1 && len(resp.ResourceList) < 1 && len(resp.SubDomain) < 1 {
		return nil, fmt.Errorf("空数据")
	}
	return &dao.WebInfo{
		PortID: portID, StartURL: resp.StartUrl, Title: resp.Title, Server: resp.Server,
		ContentType: resp.ContentType, LoginList: strings.Join(resp.LoginList, ","),
		UploadList: strings.Join(resp.UploadList, ","), SubDomain: strings.Join(resp.SubDomain, ","),
		RouteList: strings.Join(resp.RouteList, ","), ResourceList: strings.Join(resp.ResourceList, ","),
	}, nil
}
