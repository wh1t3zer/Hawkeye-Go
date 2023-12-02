package module

import (
	"context"
	"fmt"

	"github.com/e421083458/golang_common/lib"
	"github.com/wh1t3zer/Hawkeye/dao"
	pb "github.com/wh1t3zer/Hawkeye/micro/proto/grpc"
	"google.golang.org/grpc"
)

// DomainScanner ...
type DomainScanner struct {
	Domain     string
	DomainDict []string

	Conn *grpc.ClientConn
}

// NewDomainScanner ...
func NewDomainScanner(domain string, domainDict []string) (*DomainScanner, error) {
	conn, err := NewDefaultConn(lib.GetStringConf("micro.aum.name"), lib.GetStringConf("micro.consul.registry_address"))
	if err != nil {
		return nil, err
	}
	return &DomainScanner{Domain: domain, DomainDict: domainDict, Conn: conn}, nil
}

// DomainResolv 域名解析
func (dsr DomainScanner) DomainResolv() (string, error) {
	// 2.调用方法
	call := pb.NewDomainClient(dsr.Conn)
	req := &pb.ResolvRequest{Domain: dsr.Domain}
	resp, err := call.Resolv(context.Background(), req)
	if err != nil {
		fmt.Printf("Failed call method(Resolv), info: %#v\n", err)
		return "", err
	}
	return resp.Ip, nil
}

// DomainAnls ...
func (dsr DomainScanner) DomainAnls(assetID int64) (*dao.DomainInfo, error) {
	// 2.调用方法
	call := pb.NewDomainClient(dsr.Conn)
	req := &pb.AnlsRequest{Domain: dsr.Domain, DomainDict: dsr.DomainDict}
	resp, err := call.Analysis(context.Background(), req)
	if err != nil {
		fmt.Printf("Failed call method(Analysis), info: %#v\n", err)
		return nil, err
	}
	return &dao.DomainInfo{
		AssetID: assetID, Domain: dsr.Domain, SubDomainList: resp.SubdomainList, DomainServer: resp.DomainServer,
		Registrar: resp.Registrar, RegisterDate: resp.RegisterDate, NameServer: resp.NameServer, Status: resp.Status,
	}, nil
}
