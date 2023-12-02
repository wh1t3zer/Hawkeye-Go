package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye/utils"
)

// TrojanItemOutput 木马服务
type TrojanItemOutput struct {
	AssetID    int64  `json:"asset_id" form:"asset_id"`
	AssetIP    string `json:"asset_name" form:"asset_name"`
	PortID     int64  `json:"port_id" form:"Port_id"`
	PortName   string `json:"port_name" form:"port_name"`
	RealServer string `json:"real_server" form:"real_server"`
	CreateAT   string `json:"create_at" form:"create_at"`
	SpareLine  int8   `json:"line" form:"line"` // 1是容器、穿透, 2是主机,直通
}

// TrojanListOutput ...
type TrojanListOutput struct {
	Total int64               `json:"total" form:"total" comment:"总数" example:"" validate:""` //总数
	List  []*TrojanItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   //列表
}

// TrojanConnInput 木马连接会话
type TrojanConnInput struct {
	AssetID   string `json:"id" form:"id" comment:"资产ID" example:"asset_id(\\w{24})" validate:"required"`
	AssetName string `json:"name" form:"name" comment:"资产名" example:"asset_name(weblogic)" validate:"required"`
	SpareLine int8   `json:"line" form:"line" comment:"连接线路" example:"0" validate:"required"`
}

// GetValidParam 校验新增参数,绑定结构体,校验参数
func (s *TrojanConnInput) GetValidParam(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, s)
}
