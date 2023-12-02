package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye/utils"
)

// AdminSessionInfo ...
type AdminSessionInfo struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user_name"`
	LoginTime time.Time `json:"login_time"`
}

// AdminLoginInput ...
type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"管理员用户名" example:"admin" validate:"required,valid_username"` //管理员用户名
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`                   //密码
}

// BindValidParam ...
func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, param)
}

// AdminLoginOutput ...
type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""` //token
}
