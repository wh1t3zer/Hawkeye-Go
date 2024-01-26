package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// AdminInfoOutput ...
type AdminInfoOutput struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

// ChangePwdInput ...
type ChangePwdInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"` //密码
}

// BindValidParam ...
func (param *ChangePwdInput) BindValidParam(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, param)
}
