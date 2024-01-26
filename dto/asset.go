package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// AssetListOutput ...
type AssetListOutput struct {
	List  []AssetInfoOutput `json:"list" form:"list" comment:"列表"`
	Total int64             `json:"total" form:"total" comment:"总数"`
}

// AssetInfoOutput 资产自身信息输出, 用于table表
type AssetInfoOutput struct {
	ID         int64    `json:"id"`
	TaskID     int64    `json:"task_id" description:"任务ID"`
	TaskName   string   `json:"task_name" description:"任务名"`
	TaskStatus string   `json:"task_status" description:"任务状态"`
	IP         string   `json:"ip" description:"IP"`
	GPS        string   `json:"gps" description:"GPS"`
	AREA       string   `json:"area" description:"区域"`
	ISP        string   `json:"isp" description:"运营商"`
	OS         string   `json:"os" description:"操作系统"`
	Vendor     string   `json:"vendor" description:"设备"`
	CreatedAt  string   `json:"create_at" description:"添加时间"`
	PortArray  []string `json:"port_array" description:"端口列表"`
	VulCount   int8     `json:"vul_count" description:"漏洞数量"`
}

// AssetAddInput ...
type AssetAddInput struct {
	TaskID int64  `json:"task_id" form:"task_id" comment:"任务ID" validate:"required"`
	IP     string `json:"ip" form:"ip" comment:"IP地址" validate:"required"`
	GPS    string `json:"gps" form:"gps" comment:"GPS" validate:""`
	AREA   string `json:"area" form:"area" comment:"区域" validate:""`
	ISP    string `json:"isp" form:"isp" comment:"运营商" validate:""`
	OS     string `json:"os" form:"os" comment:"操作系统" validate:""`
	Vendor string `json:"vendor" form:"vendor" comment:"设备名" validate:""`
}

// GetValidParams ...
func (params *AssetAddInput) GetValidParams(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, params)
}

// AssetUpdateInput ...
type AssetUpdateInput struct {
	ID     int64  `json:"id" form:"id" comment:"资产ID" validate:"required"`
	TaskID int64  `json:"task_id" form:"task_id" comment:"任务ID" validate:""`
	IP     string `json:"ip" form:"ip" comment:"IP地址" validate:""`
	GPS    string `json:"gps" form:"gps" comment:"GPS" validate:""`
	AREA   string `json:"area" form:"area" comment:"区域" validate:""`
	ISP    string `json:"isp" form:"isp" comment:"运营商" validate:""`
	OS     string `json:"os" form:"os" comment:"操作系统" validate:""`
	Vendor string `json:"vendor" form:"vendor" comment:"设备名" validate:""`
}

// GetValidParams ...
func (params *AssetUpdateInput) GetValidParams(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, params)
}

// 聚合查询---------------------------------------------------------

// AssetInfo [-] 资产基本信息
type AssetInfo struct {
	ID        int64     `json:"id" form:"primary_key"`
	TaskID    int64     `json:"task_id" form:"task_id" description:"任务ID"`
	Target    string    `json:"target" form:"target" description:"target"`
	IP        string    `json:"ip" form:"ip" description:"IP"`
	GPS       string    `json:"gps" form:"gps" description:"GPS"`
	AREA      string    `json:"area" form:"area" description:"区域"`
	ISP       string    `json:"isp" form:"isp" description:"运营商"`
	OS        string    `json:"os" form:"os" description:"操作系统"`
	Vendor    string    `json:"vendor" form:"vendor" description:"设备"`
	CreatedAt time.Time `json:"create_at" form:"create_at" description:"添加时间"`
}

// WebTableBoxCard Web列表
type WebTableBoxCard struct {
	Title  string          `json:"title"`  // 卡片标题
	Image  string          `json:"image"`  // 卡片头像
	Type   string          `json:"type"`   // 卡片绘图类型 pie line tatle
	Series []WebInfoOutput `json:"series"` // 数据源
}

// VulTableBoxCard 漏洞列表
type VulTableBoxCard struct {
	Title  string          `json:"title"`  // 卡片标题
	Image  string          `json:"image"`  // 卡片头像
	Type   string          `json:"type"`   // 卡片绘图类型 pie line tatle
	Series []VulInfoOutput `json:"series"` // 数据源
}

// TrapInfo 蜜罐识别信息
type TrapInfo struct {
	Name   string `json:"name"`
	Verify string `json:"verify"`
}

// AssetInfoCard ...
type AssetInfoCard struct {
	Title string `json:"title"` // 卡片标题
	Image string `json:"image"` // 卡片头像
	Type  string `json:"type"`  // 卡片绘图类型 pie line tatle
	AssetInfo
}

// AssetDetailOutput [+] 关联部分数据的详情输出
type AssetDetailOutput struct {
	Trap        *TrapInfo        `json:"trap_info"`
	AssetOutput *AssetInfoCard   `json:"asset_info" form:"asset_info" description:"资产基本信息"`
	Box1        *ChartBoxCard    `json:"box1" form:"box1" description:"域名占比/服务占比"`
	Box2        *ChartBoxCard    `json:"box2" form:"box2" description:"产品占比/漏洞占比"`
	WebList     *WebTableBoxCard `json:"web_list" form:"web_list" description:"web列表"`
	VulList     *VulTableBoxCard `json:"vul_list" form:"vul_list" description:"漏洞列表"`
}
