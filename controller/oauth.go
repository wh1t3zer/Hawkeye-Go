package controller

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye/dto"
	"github.com/wh1t3zer/Hawkeye/middleware"
	"github.com/wh1t3zer/Hawkeye/utils"
)

// OauthController ...
type OauthController struct{}

// OauthRegister ...
func OauthRegister(group *gin.RouterGroup) {
	adminLogin := &OauthController{}
	group.POST("/tokens", adminLogin.Tokens)
}

// Tokens godoc
// @Summary 获取TOKEN
// @Description 获取TOKEN
// @Tags OAUTH
// @ID /oauth/toekns
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} middleware.Response{data=dto.TokensOutput} "success"
// @Router /oauth/toekns [post]
func (oauth *OauthController) Tokens(c *gin.Context) {
	params := &dto.TokensInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	out := &dto.TokensOutput{}
	middleware.ResponseSuccess(c, out)
}

// AdminLoginOut godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (oauth *OauthController) AdminLoginOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(utils.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(c, "Logout Successfully!")
}
