package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wh1t3zer/Hawkeye/dao"
	"github.com/wh1t3zer/Hawkeye/dto"
	"github.com/wh1t3zer/Hawkeye/middleware"
	"github.com/wh1t3zer/Hawkeye/module"
)

// TaskController ...
type TaskController struct{}

// TaskRegister ...
func TaskRegister(router *gin.RouterGroup) {
	ctl := TaskController{}
	router.GET("/list", ctl.TaskList)
	router.GET("/info", ctl.TaskInfo)
	router.GET("/detail", ctl.TaskDetail)
	router.POST("/add", ctl.TaskAdd)
	router.PUT("/update", ctl.TaskUpdate)
	router.DELETE("/delete", ctl.TaskDelete)
	router.GET("/stat", ctl.TaskStat)
	router.PATCH("/exec", ctl.TaskExec) // 执行任务
}

// TaskList godoc
// @Summary 任务列表
// @Description 任务列表
// @Tags 任务管理
// @ID /task/list
// @Accept  json
// @Produce  json
// @Param info query string false "模糊查询"
// @Param limit query string true "每页条数"
// @Param page query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.TaskListOutput} "success"
// @Router /task/list [get]
func (admin *TaskController) TaskList(c *gin.Context) {
	params := &dto.PublicListInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	info := &dao.TaskInfo{}
	list, total, err := info.PageList(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	taskList := []dto.TaskInfoOutput{}
	for _, item := range list {
		assetinfo := &dao.AssetInfo{TaskID: item.ID}
		_, total, err := assetinfo.AllRecord(c, lib.GORMDefaultPool)
		if err != nil {
			fmt.Println(err)
		}
		taskList = append(taskList, dto.TaskInfoOutput{
			ID: item.ID, RuleID: item.RuleID, Name: item.Name, TargetList: item.TargetList, IsDelete: item.IsDelete,
			WebScan: item.WebScan, PocScan: item.PocScan, AuthScan: item.AuthScan, TrapScan: item.TrapScan, AssetNum: int8(total),
			Recursion: item.Recursion, Progress: item.Progress, Percent: item.Percent, Status: item.Status,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"), UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	output := dto.TaskListOutput{
		List:  taskList,
		Total: total,
	}
	middleware.ResponseSuccess(c, output)
	return
}

// TaskInfo godoc
// @Summary 任务详情
// @Description 任务详情
// @Tags 任务管理
// @ID /task/info
// @Accept  json
// @Produce  json
// @Param id query string true "任务ID"
// @Success 200 {object} middleware.Response{data=dao.TaskInfo} "success"
// @Router /task/info [get]
func (admin *TaskController) TaskInfo(c *gin.Context) {
	params := &dto.PublicIDInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.TaskInfo{ID: params.ID}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, info)
	return
}

// TaskDelete godoc
// @Summary 任务删除
// @Description 任务删除
// @Tags 任务管理
// @ID /task/delete
// @Accept  json
// @Produce  json
// @Param id query string true "任务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /task/delete [delete]
func (admin *TaskController) TaskDelete(c *gin.Context) {
	params := &dto.PublicIDInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.TaskInfo{ID: params.ID}
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
	middleware.ResponseSuccess(c, "")
	return
}

// TaskAdd godoc
// @Summary 任务添加
// @Description 任务添加
// @Tags 任务管理
// @ID /task/add
// @Accept  json
// @Produce  json
// @Param body body dto.TaskAddInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /task/add [post]
func (admin *TaskController) TaskAdd(c *gin.Context) {
	params := &dto.TaskAddInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	info := &dao.TaskInfo{
		RuleID: params.RuleID, Name: params.Name, TargetList: params.TargetList, WebScan: params.WebScan, Status: "New",
		PocScan: params.PocScan, AuthScan: params.AuthScan, TrapScan: params.TrapScan, Recursion: params.Recursion,
	}
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, "success")
	return
}

// TaskUpdate godoc
// @Summary 任务更新
// @Description 任务更新
// @Tags 任务管理
// @ID /task/update
// @Accept  json
// @Produce  json
// @Param body body dto.TaskUpdateInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /task/update [put]
func (admin *TaskController) TaskUpdate(c *gin.Context) {
	params := &dto.TaskUpdateInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.TaskInfo{ID: params.ID}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	info.RuleID = params.RuleID
	info.Name = params.Name
	info.TargetList = params.TargetList
	info.WebScan = params.WebScan
	info.PocScan = params.PocScan
	info.AuthScan = params.AuthScan
	info.TrapScan = params.TrapScan
	info.Recursion = params.Recursion
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "success")
	return
}

// TaskExec godoc
// @Summary 执行任务
// @Description 执行任务
// @Tags 任务管理
// @ID /task/exec
// @Accept  json
// @Produce  json
// @Param id query string true "任务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /task/exec [patch]
func (admin *TaskController) TaskExec(c *gin.Context) {
	params := &dto.PublicIDInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 1. 找任务记录
	search := &dao.TaskInfo{ID: params.ID}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// 2. 将原有的资产进行假删除
	asset := &dao.AssetInfo{TaskID: info.ID}
	array, _, err := asset.AllRecord(c, lib.GORMDefaultPool)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	for _, item := range array {
		item.IsDelete = 1
		if err := item.Save(c, lib.GORMDefaultPool); err != nil {
			middleware.ResponseError(c, 2004, err)
		}
	}

	// 3. 更新任务状态
	info.Status = "Processing"
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2005, err)
		return
	}
	// 4. 执行任务
	executor := module.InitExecutor(info)
	fmt.Println("succe..")
	go executor.Run(c)
	// Websocket 实时传送扫描结果数据
	middleware.ResponseSuccess(c, "success")
	return
}

// TaskDetail godoc
// @Summary 任务详情(静态视图)
// @Description 任务详情(静态视图)
// @Tags 任务管理
// @ID /task/detail
// @Accept  json
// @Produce  json
// @Param id query string true "任务ID"
// @Success 200 {object} middleware.Response{data=dto.TaskDashboardOutput} "success"
// @Router /task/detail [get]
func (admin *TaskController) TaskDetail(c *gin.Context) {
	params := &dto.PublicIDInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 查库
	task := &dao.TaskInfo{ID: params.ID}
	taskinfo, err := task.Find(c, lib.GORMDefaultPool, task)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	out := taskinfo.TaskStatistics(c)
	middleware.ResponseSuccess(c, out)
	return
}

// TaskStat godoc
// @Summary 任务视图
// @Description 任务视图
// @Tags 任务管理
// @ID /task/stat/get
// @Accept  json
// @Produce  json
// @Param id query string true "任务ID"
// @Success 200 {object} middleware.Response{data=dto.TaskDashboardOutput} "success"
// @Router /task/stat [get]
func (admin *TaskController) TaskStat(c *gin.Context) {
	params := &dto.PublicIDInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 2.升级get请求为webSocket协议
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
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
		task := &dao.TaskInfo{ID: params.ID} //task会被回收
		taskinfo, _ := task.Find(c, lib.GORMDefaultPool, task)

		// 输出
		out := taskinfo.TaskStatistics(c)
		payload, _ := json.Marshal(out)
		// fmt.Printf("持续{%v, %v}更新数据,直到客户端关闭\n", params.ID, _id)
		if err = ws.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Println("Write error, info: ", err)
			return
		}
		time.Sleep(time.Millisecond * 1000)

	}
	fmt.Println("借宿了")
	return
}
