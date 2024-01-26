package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// TaskListOutput ...
type TaskListOutput struct {
	List  []TaskInfoOutput `json:"list" form:"list" comment:"列表"`
	Total int64            `json:"total" form:"total" comment:"总数"`
}

// TaskInfoOutput ...
type TaskInfoOutput struct {
	ID         int64  `json:"id"`
	RuleID     int64  `json:"rule_id" description:"规则id"`
	Name       string `json:"name" description:"任务名"`
	TargetList string `json:"target_list" description:"目标列表"`
	WebScan    int8   `json:"web_scan" description:"Web扫描"`
	PocScan    int8   `json:"poc_scan" description:"Poc扫描"`
	AuthScan   int8   `json:"auth_scan" description:"权限扫描"`
	TrapScan   int8   `json:"trap_scan" description:"蜜罐识别"`
	Recursion  int8   `json:"recursion" description:"扫描周期"`
	Progress   string `json:"progress" description:"扫描进程"`
	Percent    int8   `json:"percent" description:"扫描百分比0-100"`
	Status     string `json:"status" description:"扫描状态"`
	CreatedAt  string `json:"create_at" description:"添加时间"`
	UpdatedAt  string `json:"update_at" description:"更新时间"`
	IsDelete   int8   `json:"is_delete" description:"是否已删除；0：否；1：是"`
	AssetNum   int8   `json:"asset_num" description:"资产数量"`
}

// TaskAddInput ...
type TaskAddInput struct {
	RuleID     int64  `json:"rule_id" form:"rule_id" comment:"规则id" validate:""`
	Name       string `json:"name" form:"name" comment:"任务名" validate:"required"`
	TargetList string `json:"target_list" form:"target_list" comment:"目标列表" validate:"required"`
	WebScan    int8   `json:"web_scan" form:"web_scan" comment:"Web扫描" validate:""`
	PocScan    int8   `json:"poc_scan" form:"poc_scan" comment:"Poc扫描" validate:""`
	AuthScan   int8   `json:"auth_scan" form:"auth_scan" comment:"权限扫描" validate:""`
	TrapScan   int8   `json:"trap_scan" form:"trap_scan" comment:"蜜罐识别" validate:""`
	Recursion  int8   `json:"recursion" form:"recursion" comment:"扫描周期" validate:""`
}

// GetValidParams ...
func (params *TaskAddInput) GetValidParams(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, params)
}

// TaskUpdateInput ...
type TaskUpdateInput struct {
	ID         int64  `json:"id" form:"id" comment:"任务ID" validate:"required"`
	RuleID     int64  `json:"rule_id" form:"rule_id" comment:"web规则id" validate:""`
	Name       string `json:"name" form:"name" comment:"任务名" validate:"required"`
	TargetList string `json:"target_list" form:"target_list" comment:"目标列表" validate:"required"`
	WebScan    int8   `json:"web_scan" form:"web_scan" comment:"Web扫描" validate:""`
	PocScan    int8   `json:"poc_scan" form:"poc_scan" comment:"Poc扫描" validate:""`
	AuthScan   int8   `json:"auth_scan" form:"auth_scan" comment:"权限扫描" validate:""`
	TrapScan   int8   `json:"trap_scan" form:"trap_scan" comment:"蜜罐识别" validate:""`
	Recursion  int8   `json:"recursion" form:"recursion" comment:"扫描周期" validate:""`
}

// GetValidParams ...
func (params *TaskUpdateInput) GetValidParams(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, params)
}
