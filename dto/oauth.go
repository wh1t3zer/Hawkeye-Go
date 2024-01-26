package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// TokensInput ...
type TokensInput struct {
	GrantType string `json:"grant_type" form:"grant_type" comment:"授权类型" example:"client_credentials" validate:"required"` //授权类型
	Scope     string `json:"scope" form:"scope" comment:"权限范围" example:"read_write" validate:"required"`                   //权限范围
}

// BindValidParam ...
func (param *TokensInput) BindValidParam(c *gin.Context) error {
	return utils.DefaultGetValidParams(c, param)
}

// TokensOutput ...
type TokensOutput struct {
	AccessToken string `json:"access_token" form:"access_token"` //access_token
	ExpiresIn   int    `json:"expires_in" form:"expires_in"`     //超时时间
	TokenType   string `json:"token_type" form:"token_type"`     //类型
	Scope       string `json:"scope" form:"scope"`               //权限范围
}
