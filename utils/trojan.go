package utils

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/micro/go-micro/registry"
)

// TrojanInfoHandler 处理Torjan服务并返回信息
func TrojanInfoHandler(value *registry.Service) (portID int64, err error) {
	// fmt.Printf("%#v\n", service)
	// fmt.Printf("%v - %v - %v - %v\n", service.Name(), service.Options(), service.String(), service.Server().Options().Metadata)
	// regexp
	fmt.Println(value.Name)
	re, _ := regexp.Compile("^\\d+")
	if !re.MatchString(value.Name) {
		return 0, fmt.Errorf("not trojan service")
	}
	fmt.Println("acooo", value.Name)
	a, err := strconv.Atoi(value.Name)
	if err != nil {
		return 0, err
	}
	return int64(a), nil
}
