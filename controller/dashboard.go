package controller

// 默认是所有的数据统计
// 可选择是某个任务的数据统计
// 也可以是某个资产的数据统计

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wh1t3zer/Hawkeye/dao"
	"github.com/wh1t3zer/Hawkeye/dto"
	"github.com/wh1t3zer/Hawkeye/middleware"
)

// DashboardController ...
type DashboardController struct{}

// DashboardRegister ...
func DashboardRegister(group *gin.RouterGroup) {
	assetCtl := &DashboardController{}
	group.GET("/all", assetCtl.DefaultDashboard)
}

// DefaultDashboard godoc
// @Summary 所有数据统计(实时)
// @Description 所有数据统计(实时)
// @Tags 首页大盘
// @ID /dashboard/all
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.DashboardOutput} "success"
// @Router /dashboard/all [get]
func (assetctl *DashboardController) DefaultDashboard(c *gin.Context) {
	// 0.升级get请求为webSocket协议
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	defer ws.Close()

	broke := false
	go func() { // 监听客户端是否关闭链接
		if _, _, err := ws.ReadMessage(); err != nil {
			broke = true
		}
	}()
	for !broke {
		out := &dto.DashboardOutput{
			PanelGroup: &dto.PanelGroupData{},
			Box1: &dto.ChartBoxCard{
				Title: "Hardware", Image: "https://zan71.com/cdn-img/icon/avatar/tx1.gif", Type: "pie", Series: []dto.ChartSeries{},
			},
			Box2: &dto.ChartBoxCard{
				Title: "Software", Image: "https://zan71.com/cdn-img/icon/avatar/tx.gif", Type: "pie", Series: []dto.ChartSeries{},
			},
			Box3: &dto.ChartBoxCard{
				Title: "SubDomain", Image: "https://zan71.com/cdn-img/icon/avatar/tx1.gif", Type: "pie", Series: []dto.ChartSeries{},
			},
			Box4: &dto.ChartBoxCard{
				Title: "Vul Type", Image: "https://zan71.com/cdn-img/icon/avatar/tx1.gif", Type: "pie", Series: []dto.ChartSeries{},
			},
			Box5: &dto.TableBoxCard{
				Title: "Latest Vul", Image: "https://zan71.com/cdn-img/icon/avatar/tx1.gif", Type: "table", Series: []dto.VulInfoOutput{},
			},
			Box6: &dto.ChartBoxCard{
				Title: "Web Site", Image: "https://zan71.com/cdn-img/icon/avatar/tx.gif", Type: "line", Series: []dto.ChartSeries{},
			},
		}
		// 1.查资产
		asset := &dao.AssetInfo{}
		assetArray, atotal, err := asset.AllRecord(c, lib.GORMDefaultPool)
		if err != nil {
			break
		}

		vendors := make(map[string]int, 10)
		for _, asset := range assetArray {
			vendors[asset.Vendor]++
		}

		// 2.查域名
		domain := &dao.DomainInfo{}
		domainArray, _, err := domain.AllRecord(c, lib.GORMDefaultPool)
		if err != nil {
			break
		}
		domains := make(map[string]int, 10)
		for _, domain := range domainArray {
			for _, item := range strings.Split(domain.SubDomainList, ",") {
				if item != "" {
					domains[item]++
				}
			}
		}

		// 3.查端口表
		srv := &dao.PortInfo{}
		portArray, stotal, err := srv.AllRecord(c, lib.GORMDefaultPool)
		if err != nil {
			break
		}
		// 2.2.1 遍历端口表
		softwares := make(map[string]int, 10)
		for _, portinfo := range portArray {
			software := portinfo.Product
			if software == "" {
				software = portinfo.Name
			}
			softwares[software]++
		}

		// 4.Web信息
		web := &dao.WebInfo{}
		webArray, _, err := web.AllRecord(c, lib.GORMDefaultPool)
		if err != nil {
			break
		}
		webinfos := make(map[string]int, 10)
		webinfos["Web Server"] = 0
		webinfos["Content Type"] = 0
		webinfos["Login Page"] = 0
		webinfos["Upload Page"] = 0
		webinfos["Sub Domain"] = 0
		webinfos["Site URL"] = 0
		webinfos["Resource Path"] = 0
		for _, info := range webArray {
			if strings.TrimSpace(info.Server) != "" {
				webinfos["Web Server"] += len(strings.Split(info.Server, ","))
			}
			if strings.TrimSpace(info.ContentType) != "" {
				webinfos["Content Type"] += len(strings.Split(info.ContentType, ","))
			}
			if strings.TrimSpace(info.LoginList) != "" {
				webinfos["Login Page"] += len(strings.Split(info.LoginList, ","))
			}
			if strings.TrimSpace(info.UploadList) != "" {
				webinfos["Upload Page"] += len(strings.Split(info.UploadList, ","))
			}
			if strings.TrimSpace(info.SubDomain) != "" {
				webinfos["Sub Domain"] += len(strings.Split(info.SubDomain, ","))
			}
			if strings.TrimSpace(info.RouteList) != "" {
				webinfos["Site URL"] += len(strings.Split(info.RouteList, ","))
			}
			if strings.TrimSpace(info.ResourceList) != "" {
				webinfos["Resource Path"] += len(strings.Split(info.ResourceList, ","))
			}
		}
		// 5.找漏洞信息, 漏洞类型、漏洞列表
		vulinfo := &dao.VulInfo{}
		vulArray, vtotal, err := vulinfo.AllRecord(c, lib.GORMDefaultPool)
		vultypes := make(map[string]int, 10)
		outarray := []dto.VulInfoOutput{}
		for _, vul := range vulArray {
			// 漏洞类型
			poctObj := &dao.PocPlugin{ID: vul.PluginID}
			poc, _ := poctObj.Find(c, lib.GORMDefaultPool, poctObj)
			vultypes[poc.VulType]++
			// 漏洞列表
			assetObj := &dao.AssetInfo{ID: vul.AssetID}
			asset, _ := assetObj.Find(c, lib.GORMDefaultPool, assetObj)

			portObj := &dao.PortInfo{ID: vul.PortID}
			port, _ := portObj.Find(c, lib.GORMDefaultPool, portObj)
			outarray = append(outarray, dto.VulInfoOutput{
				ID:              vul.ID,
				AssetID:         vul.AssetID,
				Asset:           fmt.Sprintf("%v:%v", asset.IP, port.Port),
				PortID:          vul.PortID,
				PluginID:        vul.PluginID,
				AppName:         poc.AppName,
				VulName:         poc.VulName,
				VulType:         poc.VulType,
				VerifyURL:       vul.VerifyURL,
				VerifyPayload:   vul.VerifyPayload,
				VerifyResult:    vul.VerifyResult,
				ExploitURL:      vul.ExploitURL,
				ExploitPayload:  vul.ExploitPayload,
				ExploitResult:   vul.ExploitResult,
				WebshellURL:     vul.WebshellURL,
				WebshellPayload: vul.WebshellPayload,
				WebshellResult:  vul.WebshellResult,
				TrojanURL:       vul.TrojanURL,
				TrojanPayload:   vul.TrojanPayload,
				TrojanResult:    vul.TrojanResult,
				CreatedAt:       vul.CreatedAt.Format("2006-01-02 15:04:05"),
				IsDelete:        vul.IsDelete,
			})
		}

		// 组装数据
		out.PanelGroup.AssetCount = int(atotal)
		out.PanelGroup.ServiceCount = int(stotal)
		out.PanelGroup.VulCount = int(vtotal)
		out.Box5.Series = outarray
		for key, value := range vendors {
			out.Box1.Series = append(out.Box1.Series, dto.ChartSeries{Name: key, Value: value})
		}
		for key, value := range softwares {
			out.Box2.Series = append(out.Box2.Series, dto.ChartSeries{Name: key, Value: value})
		}
		for key, value := range domains {
			out.Box3.Series = append(out.Box3.Series, dto.ChartSeries{Name: key, Value: value})
		}
		for key, value := range vultypes {
			out.Box4.Series = append(out.Box4.Series, dto.ChartSeries{Name: key, Value: value})
		}
		for key, value := range webinfos {
			out.Box6.Series = append(out.Box6.Series, dto.ChartSeries{Name: key, Value: value})
		}
		payload, _ := json.Marshal(out)
		if err = ws.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Println("Write error, info: ", err)
			break
		}

		time.Sleep(time.Millisecond * 5000)
	}
	log.Println("Dashboard websocket broker!!!")
	return
}
