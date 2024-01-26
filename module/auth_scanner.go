package module

import (
	"context"
	"fmt"

	"github.com/e421083458/golang_common/lib"
	"github.com/wh1t3zer/Hawkeye-Go/dao"
	pb "github.com/wh1t3zer/Hawkeye-Go/micro/proto/grpc"
	"google.golang.org/grpc"
)

// AuthScanner ...
type AuthScanner struct {
	Args     string
	UserName []string
	Password []string
	Conn     *grpc.ClientConn
}

// NewAuthScanner ...
func NewAuthScanner() (*AuthScanner, error) {
	conn, err := NewDefaultConn(lib.GetStringConf("micro.aum.name"), lib.GetStringConf("micro.consul.registry_address"))
	if err != nil {
		return nil, err
	}
	return &AuthScanner{
		Args: "", Conn: conn, // Target: []string{fmt.Sprintf("%v:%v", host, port)}
		UserName: []string{"root", "admin", "test"}, Password: []string{"123456", "root", "admin", "test"},
	}, nil
}

// Brute 字典爆破
func (v AuthScanner) Brute(assetid, portid int64, target, service string) ([]*dao.AuthInfo, error) {
	result := []*dao.AuthInfo{}
	// 2.调用方法
	call := pb.NewVulClient(v.Conn)

	req := &pb.AuthRequest{
		Service: service, Args: v.Args, TargetList: []string{target}, UsernameList: v.UserName, PasswordList: v.Password,
	}
	resp, err := call.Hydra(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("Failed call method(Brute), info: %#v", err)
	}

	for _, item := range resp.Array {
		result = append(result, &dao.AuthInfo{
			AssetID: assetid, PortID: portid, Target: item.Target, Service: item.Service,
			Username: item.Username, Password: item.Password, Command: item.Command,
		})
	}

	return result, nil
}
