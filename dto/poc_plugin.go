package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// PocPluginAddInput ...
type PocPluginAddInput struct {
	VulID          string `json:"vul_id" form:"vul_id" comment:"vul_id" validate:"required"`
	VulName        string `json:"vul_name" form:"vul_name" comment:"漏洞名" validate:"required"`
	VulType        string `json:"vul_type" form:"vul_type" comment:"漏洞类型" validate:"required"`
	VulDate        string `json:"vul_date" form:"vul_date" comment:"漏洞发布日期" validate:""` // 2016-01-02
	Version        string `json:"version" form:"version" comment:"插件本别" validate:""`
	Author         string `json:"author" form:"author" comment:"编写者" validate:"required"`
	AppPowerLink   string `json:"app_powerLink" form:"app_powerLink" comment:"产商链接" validate:""`
	AppName        string `json:"app_name" form:"app_name" comment:"应用名" validate:"required"`
	AppVersion     string `json:"app_version" form:"app_version" comment:"应用版本" validate:"required"`
	Desc           string `json:"desc" form:"desc" comment:"漏洞描述" validate:"required"`
	Cnnvd          string `json:"cnnvd" form:"cnnvd" comment:"cnnvd" validate:""`
	CveID          string `json:"cve_id" form:"cve_id" comment:"cve_id" validate:""`
	Rank           int8   `json:"rank" form:"rank" comment:"危险等级" validate:""`
	DefaultPorts   string `json:"default_ports" form:"default_ports" comment:"默认端口" validate:"required"`
	DefaultService string `json:"default_service" form:"default_service" comment:"默认服务" validate:"required"`
	Content        string `json:"content" form:"content" comment:"脚本内容" validate:"required"`
}

// GetValidParams ...
func (params *PocPluginAddInput) GetValidParams(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, params)
}

/*-----------------------------------输出------------------------------------------*/

// PocPluginInfoOutput 对象item信息
type PocPluginInfoOutput struct {
	ID             int64  `json:"id" form:"id"` //id
	VulID          string `json:"vul_id" form:"vul_id" description:"ssvid"`
	VulName        string `json:"vul_name" form:"vul_name" description:"漏洞名"`
	VulType        string `json:"vul_type" form:"vul_type" description:"漏洞类型"`
	VulDate        string `json:"vul_date" form:"vul_date" description:"漏洞发布日期"`
	Version        string `json:"version" form:"version" description:"插件版本"`
	Author         string `json:"author" form:"author" description:"编写者"`
	AppPowerLink   string `json:"app_powerLink" form:"app_powerLink" description:"应用链接"`
	AppName        string `json:"app_name" form:"app_name" description:"应用名"`
	AppVersion     string `json:"app_version" form:"app_version" description:"应用版本"`
	Desc           string `json:"desc" form:"desc" description:"描述"`
	Cnnvd          string `json:"cnnvd" form:"cnnvd" description:"CNNVD"`
	CveID          string `json:"cve_id" form:"cve_id" description:"CVE-ID"`
	Rank           int8   `json:"rank" form:"rank" description:"威胁等级"`
	DefaultPorts   string `json:"default_ports" form:"default_ports" description:"默认端口"`
	DefaultService string `json:"default_service" form:"default_service" description:"默认服务"`
	Content        string `json:"content" form:"content" description:"插件内容"`
	UpdatedAt      string `json:"update_at" form:"update_at" description:"更新时间"`
	CreatedAt      string `json:"create_at" form:"create_at" description:"创建时间"`
	IsDelete       int8   `json:"is_delete" form:"is_delete" description:"是否删除"`
}

// PocPluginListOutput ...
type PocPluginListOutput struct {
	Total int64                 `json:"total" form:"total" comment:"总数" example:"" validate:""` //总数
	List  []PocPluginInfoOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   //列表
}
