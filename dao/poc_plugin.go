package dao

import (
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// PocPlugin ...
type PocPlugin struct {
	ID             int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	VulID          string    `json:"vul_id" gorm:"column:vul_id" description:"vul_id"`
	VulName        string    `json:"vul_name" gorm:"column:vul_name" description:"漏洞名"`
	VulType        string    `json:"vul_type" gorm:"column:vul_type" description:"漏洞类型"`
	VulDate        time.Time `json:"vul_date" gorm:"column:vul_date" description:"漏洞发布日期"`
	Version        string    `json:"version" gorm:"column:version" description:"插件本别"`
	Author         string    `json:"author" gorm:"column:author" description:"编写者"`
	AppPowerLink   string    `json:"app_powerLink" gorm:"column:app_powerLink" description:"产商链接"`
	AppName        string    `json:"app_name" gorm:"column:app_name" description:"应用名"`
	AppVersion     string    `json:"app_version" gorm:"column:app_version" description:"应用版本"`
	Desc           string    `json:"desc" gorm:"column:desc" description:"漏洞描述"`
	Cnnvd          string    `json:"cnnvd" gorm:"column:cnnvd" description:"cnnvd"`
	CveID          string    `json:"cve_id" gorm:"column:cve_id" description:"cve_id"`
	Rank           int8      `json:"rank" gorm:"column:rank" description:"危险等级"`
	DefaultPorts   string    `json:"default_ports" gorm:"column:default_ports" description:"默认端口"`
	DefaultService string    `json:"default_service" gorm:"column:default_service" description:"默认服务"`
	Content        string    `json:"content" gorm:"column:content" description:"脚本内容"`
	UpdatedAt      time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt      time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete       int8      `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

// TableName ...
func (t *PocPlugin) TableName() string {
	return "poc_plugin"
}

// Find ...
func (t *PocPlugin) Find(c *gin.Context, tx *gorm.DB, search *PocPlugin) (*PocPlugin, error) {
	out := &PocPlugin{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Save ...
func (t *PocPlugin) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error
}

// Delete ...
func (t *PocPlugin) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// AllRecord 所有记录, 除去被删除之外(接受有ID或无ID)
func (t *PocPlugin) AllRecord(c *gin.Context, tx *gorm.DB) ([]PocPlugin, int64, error) {
	var list []PocPlugin
	var count int64

	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	query = query.Where("is_delete=?", 0)
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

