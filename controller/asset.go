package controller

import (
	"fmt"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/dao"
	"github.com/wh1t3zer/Hawkeye-Go/dto"
	"github.com/wh1t3zer/Hawkeye-Go/middleware"
)

// AssetController ...
type AssetController struct{}

// AssetRegister ...
func AssetRegister(group *gin.RouterGroup) {
	ctl := &AssetController{}
	group.GET("/info", ctl.AssetInfo)     // 自身信息
	group.GET("/detail", ctl.AssetDetail) // 关联表数据的详情信息
	group.GET("/list", ctl.AssetList)     // 列表数据
	group.PUT("/update", ctl.AssetList)
	group.POST("/add", ctl.AssetList)
	group.DELETE("/delete", ctl.AssetDelete) // 删除资产
}

// AssetList godoc
// @Summary 资产列表
// @Description 资产列表
// @Tags 资产管理
// @ID /asset/list
// @Accept  json
// @Produce  json
// @Param info query string false "模糊查询"
// @Param limit query string true "每页条数"
// @Param page query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.TaskListOutput} "success"
// @Router /asset/list [get]
func (assetctl *AssetController) AssetList(c *gin.Context) {
	params := &dto.PublicListInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	info := &dao.AssetInfo{}
	list, total, err := info.PageList(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	array := []dto.AssetInfoOutput{}
	for _, item := range list {
		// 1. 找任务的id、任务名、任务状态
		taskinfo := &dao.TaskInfo{ID: item.TaskID}
		taskinfo, _ = taskinfo.Find(c, lib.GORMDefaultPool, taskinfo)
		// 2. 找端口
		templist := make([]string, 0, 10)
		portinfo := &dao.PortInfo{AssetID: item.ID}
		portlist, _, _ := portinfo.AllRecord(c, lib.GORMDefaultPool)
		var vulcount int8 = 0
		for _, val := range portlist {
			if val.Name == "" {
				templist = append(templist, fmt.Sprintf("%v", val.Port))
			} else if val.Product != "" {
				templist = append(templist, fmt.Sprintf("%v:%v", val.Port, val.Product))
			} else {
				templist = append(templist, fmt.Sprintf("%v:%v", val.Port, val.Name))
			}
			// 3.找漏洞
			vul := &dao.VulInfo{PortID: val.ID}
			_, vt, _ := vul.AllRecord(c, lib.GORMDefaultPool)
			vulcount += int8(vt)
		}
		array = append(array, dto.AssetInfoOutput{
			ID: item.ID, TaskID: item.TaskID, TaskName: taskinfo.Name, TaskStatus: taskinfo.Status,
			IP: item.IP, GPS: item.GPS, AREA: item.AREA, ISP: item.ISP, OS: item.OS, Vendor: item.Vendor,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"), PortArray: templist, VulCount: vulcount, // 有端口低危，有服务版本中危，有漏洞高危
		})
	}
	output := dto.AssetListOutput{
		List:  array,
		Total: total,
	}
	middleware.ResponseSuccess(c, output)
	return
}

// AssetInfo godoc
// @Summary 资产信息
// @Description 资产信息
// @Tags 资产管理
// @ID /asset/info
// @Accept  json
// @Produce  json
// @Param id query string true "资产ID"
// @Success 200 {object} middleware.Response{data=dto.AssetInfoOutput} "success"
// @Router /asset/info [get]
func (assetctl *AssetController) AssetInfo(c *gin.Context) {
	params := &dto.PublicIDInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.AssetInfo{ID: params.ID}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, info)
	return
}

// AssetDetail godoc
// @Summary 资产详情
// @Description 资产详情
// @Tags 资产管理
// @ID /asset/detail
// @Accept  json
// @Produce  json
// @Param id query string true "资产ID"
// @Success 200 {object} middleware.Response{data=dto.AssetDetailOutput} "success"
// @Router /asset/detail [get]
func (assetctl *AssetController) AssetDetail(c *gin.Context) {
	params := &dto.PublicIDInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	out := &dto.AssetDetailOutput{
		Trap:        &dto.TrapInfo{}, // &dto.TrapInfo{Name: trapPlugin.Name, Verify: trap.Verify},
		AssetOutput: &dto.AssetInfoCard{Title: "Asset Info", Type: "info", Image: "http://localhost:9527/img/domain.png"},
		Box1:        &dto.ChartBoxCard{Title: "Service Info", Image: "http://localhost:9527/img/service.png", Type: "pie", Series: []dto.ChartSeries{}},
		Box2:        &dto.ChartBoxCard{Title: "Vul Type", Image: "http://localhost:9527/img/vul.png", Type: "pie", Series: []dto.ChartSeries{}},
		WebList:     &dto.WebTableBoxCard{Title: "Web info", Image: "http://localhost:9527/img/domain.png", Type: "web", Series: []dto.WebInfoOutput{}},
		VulList:     &dto.VulTableBoxCard{Title: "Vulnerability", Image: "http://localhost:9527/img/hack.png", Type: "table", Series: []dto.VulInfoOutput{}},
	}

	// 0.先找主体信息
	search := &dao.AssetInfo{ID: params.ID}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	detail := dto.AssetInfo{
		ID: info.ID, TaskID: info.TaskID, Target: info.IP, IP: info.IP, GPS: info.GPS, ISP: info.ISP,
		AREA: info.AREA, OS: info.OS, Vendor: info.Vendor, CreatedAt: info.CreatedAt,
	}

	// 1.找域名信息
	domain := &dao.DomainInfo{AssetID: params.ID}
	if domain, err = domain.Find(c, lib.GORMDefaultPool, domain); err == nil {
		detail.Target = domain.Domain
	}

	// 2.找任务信息
	task := &dao.TaskInfo{ID: info.TaskID}
	task, err = task.Find(c, lib.GORMDefaultPool, task)
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}

	// 3.找蜜罐信息 -> 只需找到一个即可, 因为只需判断一个即可认定蜜罐服务
	trap := &dao.TrapInfo{AssetID: info.ID}
	trap, err = trap.Find(c, lib.GORMDefaultPool, trap)
	if err == nil {
		trapPlugin := &dao.TrapPluginInfo{ID: trap.PluginID}
		trapPlugin, err = trapPlugin.Find(c, lib.GORMDefaultPool, trapPlugin)
		if err != nil {
			middleware.ResponseError(c, 2006, err)
			return
		}
		out.Trap = &dto.TrapInfo{Name: trapPlugin.Name, Verify: trap.Verify}
	}

	// 4. 找端口信息
	srv := &dao.PortInfo{AssetID: info.ID}
	srvList, _, err := srv.AllRecord(c, lib.GORMDefaultPool)
	if err != nil {
		middleware.ResponseError(c, 2007, err)
		return
	}
	var softwares = make(map[string]int, 10)
	var outarray = []dto.VulInfoOutput{}
	var vultypes = make(map[string]int, 10)
	var webarray = []dto.WebInfoOutput{}
	for _, service := range srvList {
		// 4.1 饼图
		software := service.Product
		if software == "" {
			software = service.Name
		}
		softwares[software]++
		// 4.2 根据端口信息找web列表,找一个即可
		web := &dao.WebInfo{PortID: service.ID}
		if web, err := web.Find(c, lib.GORMDefaultPool, web); err == nil {
			// 找到就进行追加
			webarray = append(webarray, dto.WebInfoOutput{
				ID:           web.ID,
				PortID:       web.PortID,
				StartURL:     web.StartURL,
				Title:        web.Title,
				Server:       web.Server,
				ContentType:  web.ContentType,
				LoginList:    web.LoginList,
				UploadList:   web.UploadList,
				SubDomain:    web.SubDomain,
				RouteList:    web.RouteList,
				ResourceList: web.ResourceList,
			})
		}
		// 4.3 漏洞列表
		vul := &dao.VulInfo{PortID: service.ID}
		if vulArray, _, err := vul.AllRecord(c, lib.GORMDefaultPool); err == nil {
			for _, info := range vulArray {
				// 漏洞类型
				poctObj := &dao.PocPlugin{ID: info.PluginID}
				poc, _ := poctObj.Find(c, lib.GORMDefaultPool, poctObj)
				vultypes[poc.VulType]++
				// 漏洞列表
				outarray = append(outarray, dto.VulInfoOutput{
					ID:              info.ID,
					AssetID:         info.AssetID,
					Asset:           fmt.Sprintf("%v:%v", detail.IP, service.Port),
					PortID:          info.PortID,
					PluginID:        info.PluginID,
					AppName:         poc.AppName,
					VulName:         poc.VulName,
					VulType:         poc.VulType,
					VerifyURL:       info.VerifyURL,
					VerifyPayload:   info.VerifyPayload,
					VerifyResult:    info.VerifyResult,
					ExploitURL:      info.ExploitURL,
					ExploitPayload:  info.ExploitPayload,
					ExploitResult:   info.ExploitResult,
					WebshellURL:     info.WebshellURL,
					WebshellPayload: info.WebshellPayload,
					WebshellResult:  info.WebshellResult,
					TrojanURL:       info.TrojanURL,
					TrojanPayload:   info.TrojanPayload,
					TrojanResult:    info.TrojanResult,
					CreatedAt:       info.CreatedAt.Format("2006-01-02 15:04:05"),
					IsDelete:        info.IsDelete,
				})
			}
		}
	}

	// [-] 输出
	out.AssetOutput.AssetInfo = detail
	for key, val := range softwares {
		out.Box1.Series = append(out.Box1.Series, dto.ChartSeries{Name: key, Value: val})
	}
	for key, val := range vultypes {
		out.Box2.Series = append(out.Box2.Series, dto.ChartSeries{Name: key, Value: val})
	}
	out.VulList.Series = outarray
	out.WebList.Series = webarray
	middleware.ResponseSuccess(c, out)
	return
}

// AssetDelete godoc
// @Summary 资产删除
// @Description 资产删除
// @Tags 资产管理
// @ID /asset/delete
// @Accept  json
// @Produce  json
// @Param id query string true "资产ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /asset/delete [delete]
func (assetctl *AssetController) AssetDelete(c *gin.Context) {
	params := &dto.PublicIDInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 1. 找记录
	search := &dao.AssetInfo{ID: params.ID}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// 2. 删记录
	if err := info.Delete(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// AssetAdd godoc
// @Summary 资产添加
// @Description 资产添加
// @Tags 资产管理
// @ID /asset/add
// @Accept  json
// @Produce  json
// @Param body body dto.AssetAddInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /asset/add [post]
func (assetctl *AssetController) AssetAdd(c *gin.Context) {
	params := &dto.AssetAddInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	info := &dao.AssetInfo{
		TaskID: params.TaskID, IP: params.IP, GPS: params.GPS, ISP: params.ISP,
		AREA: params.AREA, OS: params.OS, Vendor: params.Vendor,
	}
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, "success")
	return
}

// AssetUpdate godoc
// @Summary 资产更新
// @Description 资产更新
// @Tags 资产管理
// @ID /asset/update
// @Accept  json
// @Produce  json
// @Param body body dto.AssetUpdateInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /asset/update [put]
func (assetctl *AssetController) AssetUpdate(c *gin.Context) {
	params := &dto.AssetUpdateInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.AssetInfo{ID: params.ID}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	info.TaskID = params.TaskID
	info.IP = params.IP
	info.GPS = params.GPS
	info.ISP = params.ISP
	info.AREA = params.AREA
	info.OS = params.OS
	info.Vendor = params.Vendor
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "success")
	return
}
