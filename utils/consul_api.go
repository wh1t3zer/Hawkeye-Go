package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/e421083458/golang_common/lib"
)

// IPv4 ...
type IPv4 struct {
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

// LanWan ...
type LanWan struct {
	Lan IPv4 `json:"lan_ipv4"`
	Wan IPv4 `json:"wan_ipv4"`
}

// Service ...
type Service struct {
	Service         string   `json:"Service"`
	Tags            []string `json:"Tags"`
	Address         string   `json:"Address"`
	TaggedAddresses LanWan   `json:"TaggedAddresses"`
}

// Result ...
type Result struct {
	Service Service `json:"Service"`
}

// ConsulAPI ...
type ConsulAPI struct {
	URL string
}

// NewConsulAPI ...
func NewConsulAPI(portID int64) *ConsulAPI {
	return &ConsulAPI{
		URL: fmt.Sprintf("http://%s/v1/health/service/%v?dc=dc1", lib.GetStringConf("micro.consul.registry_address"), portID),
	}
}

// GetRealServer ...
func (api *ConsulAPI) GetRealServer() (Address string, Port int, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", api.URL, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	text := string(body)
	re := regexp.MustCompile("^\\[|\\]$")
	out := re.ReplaceAllString(strings.TrimSpace(text), "")
	out = strings.TrimSpace(out)
	var res Result
	if err = json.Unmarshal([]byte(out), &res); err != nil {
		return
	}
	return res.Service.TaggedAddresses.Lan.Address, res.Service.TaggedAddresses.Lan.Port, nil
}
