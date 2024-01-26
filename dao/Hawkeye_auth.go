package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/dto"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// AuthInfo ...
type AuthInfo struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	AssetID  int64  `json:"asset_id" gorm:"column:asset_id" description:"资产id"`
	PortID   int64  `json:"port_id" gorm:"column:port_id" description:"端口id"`
	Target   string `json:"target" gorm:"column:target" description:"目标"`
	Service  string `json:"service" gorm:"column:service" description:"服务"`
	Username string `json:"username" gorm:"column:username" description:"用户名"`
	Password string `json:"password" gorm:"column:password" description:"密码"`
	Command  string `json:"command" gorm:"column:command" description:"验证命令"`
}

// TableName ...
func (t *AuthInfo) TableName() string {
	return "Hawkeye_auth"
}

// Find ...
func (t *AuthInfo) Find(c *gin.Context, tx *gorm.DB, search *AuthInfo) (*AuthInfo, error) {
	model := &AuthInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *AuthInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete ...
func (t *AuthInfo) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// PageList ...
func (t *AuthInfo) PageList(c *gin.Context, tx *gorm.DB, params *dto.PublicListInput) ([]AuthInfo, int64, error) {
	var list []AuthInfo
	var count int64
	pageNo := params.Page
	pageSize := params.Limit

	offset := (pageNo - 1) * pageSize
	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	if params.Info != "" {
		query = query.Where(" (ip like ? or vendor like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
	}
	err := query.Limit(pageSize).Offset(offset).Order("id asc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// AllRecord ...
func (t *AuthInfo) AllRecord(condition string, c *gin.Context, tx *gorm.DB) ([]AuthInfo, int64, error) {
	// condition为where的查询命令 port_id=1 AND asset_id=2
	var list []AuthInfo
	var count int64

	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	if condition != "" {
		// query = query.Where(fmt.Sprintf("port_id=%v", t.PortID))
		query = query.Where(condition)
	}
	err := query.Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}

	return list, count, nil
}
