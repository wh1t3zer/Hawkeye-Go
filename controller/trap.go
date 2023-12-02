package controller

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye/dao"
	"github.com/wh1t3zer/Hawkeye/dto"
	"github.com/wh1t3zer/Hawkeye/middleware"
)

// TrapController 蜜罐识别插件
type TrapController struct{}

// TrapRegister ...
func TrapRegister(router *gin.RouterGroup) {
	ctl := TrapController{}
	// 插件
	router.GET("/plugin/list", ctl.PluginList)
	router.GET("/plugin/info", ctl.PluginInfo)
	router.POST("/plugin/add", ctl.PluginAdd)
	router.PUT("/plugin/update", ctl.PluginUpdate)
	router.DELETE("/plugin/delete", ctl.PluginDelete)
}

// PluginList godoc
// @Summary 插件列表
// @Description 插件列表
// @Tags 蜜罐识别
// @ID /trap/plugin/list
// @Accept  json
// @Produce  json
// @Param info query string false "模糊查询"
// @Param limit query string true "每页条数"
// @Param page query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.TrapPluginListOutput} "success"
// @Router /trap/plugin/list [get]
func (admin *TrapController) PluginList(c *gin.Context) {
	params := &dto.PublicListInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	info := &dao.TrapPluginInfo{}
	list, total, err := info.PageList(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	pluginList := []*dto.TrapPluginOutput{}
	for _, item := range list {
		pluginList = append(pluginList, &dto.TrapPluginOutput{
			ID: item.ID, TrapID: item.TrapID, Name: item.Name, Author: item.Author, Protocol: item.Protocol,
			AppName: item.AppName, Honeypot: item.Honeypot, Desc: item.Desc, Content: item.Content, IsDelete: item.IsDelete,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"), UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	output := dto.TrapPluginListOutput{
		List:  pluginList,
		Total: total,
	}
	middleware.ResponseSuccess(c, output)
	return
}

// PluginInfo godoc
// @Summary 插件信息
// @Description 插件信息
// @Tags 蜜罐识别
// @ID /trap/plugin/info
// @Accept  json
// @Produce  json
// @Param id query string true "任务ID"
// @Success 200 {object} middleware.Response{data=dao.TaskInfo} "success"
// @Router /trap/plugin/info [get]
func (admin *TrapController) PluginInfo(c *gin.Context) {
	params := &dto.PublicIDInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.TrapPluginInfo{ID: params.ID}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, info)
	return
}

// PluginDelete godoc
// @Summary 插件删除
// @Description 插件删除
// @Tags 蜜罐识别
// @ID /trap/plugin/delete
// @Accept  json
// @Produce  json
// @Param id query string true "任务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /trap/plugin/delete [delete]
func (admin *TrapController) PluginDelete(c *gin.Context) {
	params := &dto.PublicIDInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.TrapPluginInfo{ID: params.ID}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	info.IsDelete = 1
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "delete success")
	return
}

// PluginAdd godoc
// @Summary 插件添加
// @Description 插件添加
// @Tags 蜜罐识别
// @ID /trap/plugin/add
// @Accept  json
// @Produce  json
// @Param body body dto.TrapPluginAddInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /trap/plugin/add [post]
func (admin *TrapController) PluginAdd(c *gin.Context) {
	params := &dto.TrapPluginAddInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	info := &dao.TrapPluginInfo{
		TrapID: params.TrapID, Name: params.Name, Author: params.Author, Protocol: params.Protocol,
		AppName: params.AppName, Honeypot: params.Honeypot, Desc: params.Desc, Content: params.Content,
	}
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, "success")
	return
}

// PluginUpdate godoc
// @Summary 插件更新
// @Description 插件更新
// @Tags 蜜罐识别
// @ID /trap/plugin/update
// @Accept  json
// @Produce  json
// @Param body body dto.TrapPluginUpdateInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /trap/plugin/update [put]
func (admin *TrapController) PluginUpdate(c *gin.Context) {
	params := &dto.TrapPluginUpdateInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.TrapPluginInfo{ID: params.ID}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	info.TrapID = params.TrapID
	info.Name = params.Name
	info.Author = params.Author
	info.Protocol = params.Protocol
	info.AppName = params.AppName
	info.Honeypot = params.Honeypot
	info.Desc = params.Desc
	info.Content = params.Content
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "success")
	return
}
