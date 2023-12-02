package module

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye/dao"
)

// Executor ...
type Executor struct {
	TaskID     int64
	RuleID     int64
	Name       string
	TargetList string
	WebScan    int8
	PocScan    int8
	AuthScan   int8
	TrapScan   int8
	Recursion  int8
	Progress   string
	Percent    int8
	Status     string

	step       chan string
	ccap       int8 // 基数
	history    int8
	IPPattern  *regexp.Regexp
	NetPattern *regexp.Regexp
	DomainDict []string // 域名字典查数据库
	PortDict   []string //端口字典查数据库
	AssetList  []*dao.AssetInfo
}

// InitExecutor 初始化执行器
func InitExecutor(taskinfo *dao.TaskInfo) *Executor {
	ipPattern, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")    // IP 172.31.50.249
	netPattern, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.0/\\d{1,2}") // 网段 172.31.50.0/24
	domainDict := []string{"www", "ftp", "imap", "mail", "smtp", "cdn", "cloud", "account"}
	portDict := []string{"22", "80", "443", "3306", "7001", "8081", "11211"}

	ccap := taskinfo.WebScan + taskinfo.AuthScan + taskinfo.PocScan + taskinfo.TrapScan
	stepChan := make(chan string, 10)

	return &Executor{
		TaskID: taskinfo.ID, RuleID: taskinfo.RuleID, Name: taskinfo.Name, TargetList: taskinfo.TargetList, WebScan: taskinfo.WebScan,
		PocScan: taskinfo.PocScan, AuthScan: taskinfo.AuthScan, TrapScan: taskinfo.TrapScan, Recursion: taskinfo.Recursion, ccap: ccap,
		IPPattern: ipPattern, NetPattern: netPattern, DomainDict: domainDict, PortDict: portDict, step: stepChan, Status: taskinfo.Status,
	}
}

// TargetFilter 目标过滤, 不能多网段
func (exec *Executor) TargetFilter(c *gin.Context) {
	var wg sync.WaitGroup
	// 如果是网段(todo) nmap discovery
	if netArray := exec.NetPattern.Find([]byte(exec.TargetList)); netArray != nil {
		fmt.Println("todo")
	}
	// 否则就是单ip、单域名、混合列表
	targetArray := strings.Split(exec.TargetList, ",")
	for _, target := range targetArray {
		if flag := exec.IPPattern.Find([]byte(target)); flag != nil { // 是IP
			// 1.解析主机、端口、IP信息存资产sql表
			wg.Add(1)
			go func(target string) {
				exec.IPHandler(target, c)
				wg.Done()
			}(target)
			continue
		}
		if true { // 是域名
			wg.Add(1)
			go func(target string) {
				exec.DomainHnadler(target, c)
				wg.Done()
			}(target)
			continue
		}
	}
	wg.Wait()
}

// Run ...
func (exec *Executor) Run(c *gin.Context) {
	exec.step <- "初始化任务"
	defer func() {
		exec.step <- "Successfully"
	}()

	// 0. 扫描进度协程
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300) // 任务执行最长时间5分钟，后续查规则表获得
	defer cancel()
	go exec.UpdateStatus(ctx, c)

	// 1. 目标过滤得到资产信息模型数组
	exec.TargetFilter(c)

	// 2. 资产并发扫描
	var wg sync.WaitGroup
	var host string
	for _, asset := range exec.AssetList {
		// 1.0 先找域名表、再找端口表进行组合扫描
		search1 := &dao.DomainInfo{AssetID: asset.ID}
		if domain, err := search1.Find(c, lib.GORMDefaultPool, search1); err == nil && strings.TrimSpace(domain.Domain) != "" {
			host = domain.Domain
		} else {
			host = asset.IP
		}
		fmt.Println("看看是否使用域名", host)
		search2 := &dao.PortInfo{AssetID: asset.ID}
		array, total, err := search2.AllRecord(c, lib.GORMDefaultPool)
		if total <= 0 || err != nil {
			continue
		}
		wg.Add(1)
		go func(asset *dao.AssetInfo) {
			// 1.1 执行Web扫描
			if exec.WebScan > 0 {
				exec.step <- "Web扫描Starting"
				exec.WebScanHandler(array, host, c)
				exec.step <- "Web扫描Finished"
			}
			// 1.2 执行Poc扫描
			if exec.PocScan > 0 {
				exec.step <- "漏洞扫描Starting"
				exec.PocScanHandler(asset.ID, array, host, c)
				exec.step <- "漏洞扫描Finished"
			}
			// 1.3 执行蜜罐识别
			// if exec.TrapScan > 0 {
			// 	exec.step <- "蜜罐识别Starting"
			// 	exec.TrapScanHandler(asset.ID, array, host, c)
			// 	exec.step <- "蜜罐识别Finished"
			// }
			// // 1.4 执行权限扫描
			// if exec.AuthScan > 0 {
			// 	exec.step <- "权限扫描Starting"
			// 	exec.AuthScanHandler(asset.ID, array, host, c)
			// 	exec.step <- "权限扫描Finished"
			// }
			wg.Done()
		}(asset)

	}
	wg.Wait()
}

// IPHandler 如果是IP，存资产表、存端口表
func (exec *Executor) IPHandler(ip string, c *gin.Context) {
	scanner, err := NewAssetScanner(ip, exec.PortDict)
	if err != nil {
		return
	}
	defer scanner.Conn.Close()

	// 地理信息
	ipinfo, err := scanner.IPLocation()
	if err != nil {
		return
	}

	// 主机、端口信息
	hostinfo, portarray, err := scanner.HostDetail()
	if err != nil {
		return
	}

	result := &dao.AssetInfo{TaskID: exec.TaskID, IP: ip, AREA: ipinfo.AREA, ISP: ipinfo.ISP, GPS: ipinfo.GPS, OS: hostinfo.OS, Vendor: hostinfo.Vendor}
	if err = result.Save(c, lib.GORMDefaultPool); err != nil {
		log.Println("Failed Save AssetInfo at DB, Info:", err)
		return
	}
	// 存端口表出错，不影响资产
	for _, portinfo := range portarray {
		portinfo.AssetID = result.ID
		if err = portinfo.Save(c, lib.GORMDefaultPool); err != nil {
			log.Println("Failed Save PortInfo at DB, Info:", err)
			break
		}
	}
	exec.AssetList = append(exec.AssetList, result)
}

// DomainHnadler 如果是域名，存域名表、存资产表
func (exec *Executor) DomainHnadler(domain string, c *gin.Context) {
	scanner1, err := NewDomainScanner(domain, exec.DomainDict)
	if err != nil {
		return
	}
	defer scanner1.Conn.Close()

	// 1、域名解析
	ip, err := scanner1.DomainResolv()
	if err != nil {
		return
	}
	scanner2, err := NewAssetScanner(ip, exec.PortDict)
	if err != nil {
		return
	}
	defer scanner2.Conn.Close()

	// 地理信息
	ipinfo, err := scanner2.IPLocation()
	if err != nil {
		return
	}
	// 主机、端口信息
	hostinfo, portarray, err := scanner2.HostDetail()
	if err != nil {
		return
	}

	result := &dao.AssetInfo{TaskID: exec.TaskID, IP: ip, GPS: ipinfo.GPS, AREA: ipinfo.AREA, ISP: ipinfo.ISP, OS: hostinfo.OS, Vendor: hostinfo.Vendor}
	if err = result.Save(c, lib.GORMDefaultPool); err != nil {
		log.Println("Failed Save AssetInfo at DB, Info:", err)
		return
	}
	// 存端口表出错，不影响资产
	for _, portinfo := range portarray {
		portinfo.AssetID = result.ID
		if err = portinfo.Save(c, lib.GORMDefaultPool); err != nil {
			log.Println("Failed Save PortInfo at DB, Info:", err)
			break
		}
	}

	// 2、域名查询、爆破
	domainInfo, err := scanner1.DomainAnls(result.ID)
	if err != nil {
		log.Println("Failed call DomainAnls method, info: ", err)
	}
	if err = domainInfo.Save(c, lib.GORMDefaultPool); err != nil {
		log.Println("Failed Save DomainInfo, info: ", err)
	}
	exec.AssetList = append(exec.AssetList, result)
}

// WebScanHandler Web扫描, 并存web信息表
func (exec *Executor) WebScanHandler(array []dao.PortInfo, host string, c *gin.Context) {
	scanner, err := NewWebScanner()
	if err != nil {
		log.Printf("Failed create NewWebScanner, info: %v\n", err)
		return
	}
	defer scanner.Conn.Close()

	// var wg sync.WaitGroup
	for _, portinfo := range array {
		port, _ := strconv.Atoi(strings.TrimSpace(portinfo.Port))
		info, err := scanner.WebSpider(portinfo.ID, host, int32(port))
		if err != nil {
			log.Printf("Failed exec webspider, info: %v\n", err)
			return
		}
		if err := info.Save(c, lib.GORMDefaultPool); err != nil {
			log.Printf("Failed save webspider to db, info: %v\n", err)
			return
		}
	}
}

// PocScanHandler Poc漏洞扫描
func (exec *Executor) PocScanHandler(assetid int64, array []dao.PortInfo, host string, c *gin.Context) {
	scanner, err := NewVulScanner()
	if err != nil {
		log.Printf("Failed create NewVulScanner, info: %v\n", err)
		return
	}
	defer scanner.Conn.Close()

	// 1. 查poc插件
	poc := &dao.PocPlugin{}
	pocArray, total, err := poc.AllRecord(c, lib.GORMDefaultPool)
	if err != nil || total < 1 {
		log.Printf("Failed get PocPlugin, info: %v\n", err)
		return
	}
	for _, portinfo := range array {
		target := fmt.Sprintf("%v:%v", host, portinfo.Port)
		result, err := scanner.Pocsuite(assetid, portinfo.ID, target, pocArray)
		if err != nil {
			log.Printf("Failed exec Pocsuite, info: %v\n", err)
			return
		}
		fmt.Printf("poc %v-%v result: %#v\n", target, len(result), result)
		for _, vulinfo := range result {
			if err := vulinfo.Save(c, lib.GORMDefaultPool); err != nil {
				log.Printf("Failed Save Vulinfo To DB, info: %v\n", err)
				return
			}
		}
	}
}

// AuthScanHandler ...
func (exec *Executor) AuthScanHandler(assetid int64, array []dao.PortInfo, host string, c *gin.Context) {
	scanner, err := NewAuthScanner() // 后续需要传自定义规则喔
	if err != nil {
		log.Printf("Failed create NewAuthScanner, info: %v\n", err)
		return
	}
	defer scanner.Conn.Close()

	// 从sql查用户定义的协议
	protocols := []string{"ssh", "vnc", "ftp", "redis", "mysql"}
	// [todo] 需要更改远程爆破服务的协议数量-> 支持协议数组多线程
	for _, portinfo := range array {
		target := fmt.Sprintf("%v:%v", host, portinfo.Port)

		// 1.协议是用户定义的数组里
		for _, protocol := range protocols {
			if portinfo.Name != protocol {
				continue
			}
			resp, err := scanner.Brute(assetid, portinfo.ID, target, protocol)
			fmt.Printf("爆破项, len=%v, %v\n", len(resp), protocol)
			if err != nil {
				log.Printf("Failed exec Brute, info: %v\n", err)
				return
			}
			// 存数据库
			for _, info := range resp {
				if err = info.Save(c, lib.GORMDefaultPool); err != nil {
					log.Printf("Failed Save Hawkeye_auth to DB, info: %v\n", err)
					return
				}
			}
		}
	}
}

// TrapScanHandler ...
func (exec *Executor) TrapScanHandler(assetid int64, array []dao.PortInfo, host string, c *gin.Context) {
	scanner, err := NewTrapScanner()
	if err != nil {
		log.Printf("Failed create NewTrapScanner, info: %v\n", err)
		return
	}
	defer scanner.Conn.Close()

	// 1.查插件
	search := &dao.TrapPluginInfo{}
	pluginArray, _, err := search.AllRecord(c, lib.GORMDefaultPool)
	if err != nil {
		log.Printf("Failed Query TrapPlugin from DB, info: %v\n", err)
		return
	}
	// 2.执行扫描
	for _, portinfo := range array {
		target := fmt.Sprintf("%v:%v", host, portinfo.Port)
		for _, plugin := range pluginArray {
			resp, err := scanner.Run(assetid, portinfo.ID, target, &plugin)
			if err != nil || len(resp) < 1 {
				continue
			}
			for _, info := range resp {
				fmt.Printf("Find honeypot:%#v\n", resp)
				info.Save(c, lib.GORMDefaultPool)
			}
		}
	}
}

// UpdateStatus ...
func (exec *Executor) UpdateStatus(context context.Context, c *gin.Context) {
	broke := false
	for !broke {
		select {
		case progress := <-exec.step:
			fmt.Println("加载进度: ", progress)
			search := &dao.TaskInfo{ID: exec.TaskID}
			taskInfo, err := search.Find(c, lib.GORMDefaultPool, search)
			if err != nil {
				log.Println("Failed update Taskinfo, info: ", err)
				return
			}

			taskInfo.Status, taskInfo.Progress, taskInfo.Percent = exec.calculation(progress)
			if err := taskInfo.Save(c, lib.GORMDefaultPool); err != nil {
				log.Println("Failed update Taskinfo, info: ", err)
				return
			}

			if taskInfo.Status == "Successfully" {
				broke = true
			}
		default:
			time.Sleep(time.Millisecond * 200) // 200毫秒
		}
	}
}

// Progress 计算扫描进度
func (exec *Executor) calculation(progress string) (string, string, int8) {
	exec.history++
	if progress == "Successfully" {
		return "Successfully", progress, 100
	}
	num := int(exec.history) * 100 / (len(strings.Split(exec.TargetList, ",")) * int(exec.ccap))
	if num >= 90 {
		num = 90
	}
	fmt.Println("current num:", num)

	if progress == "Processing" || progress == "Stop" || progress == "Failed" {
		exec.Status = progress
		return progress, progress, int8(num)
	}
	return exec.Status, progress, int8(num)
}
