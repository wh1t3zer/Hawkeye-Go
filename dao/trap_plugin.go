package dao

import (
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye/dto"
	"github.com/wh1t3zer/Hawkeye/utils"
)

// TrapPluginInfo ...
type TrapPluginInfo struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	TrapID    string    `json:"trap_id" gorm:"column:trap_id" description:"蜜罐ID"`
	Name      string    `json:"name" gorm:"column:name" description:"插件名"`
	Author    string    `json:"author" gorm:"column:author" description:"编写者"`
	Protocol  string    `json:"protocol" gorm:"column:protocol" description:"协议"`
	AppName   string    `json:"app_name" gorm:"column:app_name" description:"应用名"`
	Honeypot  string    `json:"honeypot" gorm:"column:honeypot" description:"蜜罐名"`
	Desc      string    `json:"desc" gorm:"column:desc" description:"描述"`
	Content   string    `json:"content" gorm:"column:content" description:"脚本内容"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

// TableName ...
func (t *TrapPluginInfo) TableName() string {
	return "trap_plugin"
}

// Find ...
func (t *TrapPluginInfo) Find(c *gin.Context, tx *gorm.DB, search *TrapPluginInfo) (*TrapPluginInfo, error) {
	model := &TrapPluginInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *TrapPluginInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete ...
func (t *TrapPluginInfo) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// PageList ...
func (t *TrapPluginInfo) PageList(c *gin.Context, tx *gorm.DB, params *dto.PublicListInput) ([]TrapPluginInfo, int64, error) {
	var list []TrapPluginInfo
	var count int64
	pageNo := params.Page
	pageSize := params.Limit

	offset := (pageNo - 1) * pageSize
	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	if params.Info != "" {
		query = query.Where(" (id like ? or asset_id like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
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
func (t *TrapPluginInfo) AllRecord(c *gin.Context, tx *gorm.DB) ([]TrapPluginInfo, int64, error) {
	var list []TrapPluginInfo
	var count int64

	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")

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
