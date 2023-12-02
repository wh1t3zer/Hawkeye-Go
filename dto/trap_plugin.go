package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye/utils"
)

// TrapPluginOutput 木马服务
type TrapPluginOutput struct {
	ID        int64  `json:"id" form:"id"`
	TrapID    string `json:"trap_id" form:"trap_id" description:"蜜罐ID"`
	Name      string `json:"name" form:"name" description:"插件名"`
	Author    string `json:"author" form:"author" description:"编写者"`
	Protocol  string `json:"protocol" form:"protocol" description:"协议"`
	AppName   string `json:"app_name" form:"app_name" description:"应用名"`
	Honeypot  string `json:"honeypot" form:"honeypot" description:"蜜罐名"`
	Desc      string `json:"desc" form:"desc" description:"描述"`
	Content   string `json:"content" form:"content" description:"脚本内容"`
	CreatedAt string `json:"create_at" form:"create_at" description:"添加时间"`
	UpdatedAt string `json:"update_at" form:"update_at" description:"更新时间"`
	IsDelete  int8   `json:"is_delete" form:"is_delete" description:"是否已删除；0：否；1：是"`
}

// TrapPluginListOutput ...
type TrapPluginListOutput struct {
	Total int64               `json:"total" form:"total" comment:"总数" example:"" validate:""` //总数
	List  []*TrapPluginOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   //列表
}

// TrapPluginAddInput ...
type TrapPluginAddInput struct {
	TrapID   string `json:"trap_id" form:"trap_id" comment:"蜜罐ID" validate:"required"`
	Name     string `json:"name" form:"name" comment:"插件名" validate:"required"`
	Author   string `json:"author" form:"author" comment:"编写者" validate:"required"`
	Protocol string `json:"protocol" form:"protocol" comment:"协议" validate:"required"`
	AppName  string `json:"app_name" form:"app_name" comment:"应用名" validate:"required"`
	Honeypot string `json:"honeypot" form:"honeypot" comment:"蜜罐名" validate:"required"`
	Desc     string `json:"desc" form:"desc" comment:"描述" validate:"required"`
	Content  string `json:"content" form:"content" comment:"脚本内容" validate:"required"`
}

// GetValidParams ...
func (params *TrapPluginAddInput) GetValidParams(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, params)
}

// TrapPluginUpdateInput ...
type TrapPluginUpdateInput struct {
	ID       int64  `json:"id" form:"id" comment:"插件ID" validate:"required"`
	TrapID   string `json:"trap_id" form:"trap_id" comment:"蜜罐ID" validate:"required"`
	Name     string `json:"name" form:"name" comment:"插件名" validate:"required"`
	Author   string `json:"author" form:"author" comment:"编写者" validate:"required"`
	Protocol string `json:"protocol" form:"protocol" comment:"协议" validate:"required"`
	AppName  string `json:"app_name" form:"app_name" comment:"应用名" validate:"required"`
	Honeypot string `json:"honeypot" form:"honeypot" comment:"蜜罐名" validate:"required"`
	Desc     string `json:"desc" form:"desc" comment:"描述" validate:"required"`
	Content  string `json:"content" form:"content" comment:"脚本内容" validate:"required"`
}

// GetValidParams ...
func (params *TrapPluginUpdateInput) GetValidParams(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, params)
}
