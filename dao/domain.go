package dao

import (
	"fmt"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// DomainInfo ...
type DomainInfo struct {
	ID            int64  `json:"id" gorm:"primary_key"`
	AssetID       int64  `json:"asset_id" gorm:"column:asset_id" description:"资产ID"`
	Domain        string `json:"domain" gorm:"column:domain" description:"域名"`
	SubDomainList string `json:"subdomain_list" gorm:"column:subdomain_list" description:"子域列表"`
	Registrar     string `json:"registrar" gorm:"column:registrar" description:"注册商"`
	RegisterDate  string `json:"register_date" gorm:"column:register_date" description:"注册日期"`
	NameServer    string `json:"name_server" gorm:"column:name_server" description:"DNS解析地址"`
	DomainServer  string `json:"domain_server" gorm:"column:domain_server" description:"域名解析器"`
	Status        string `json:"status" gorm:"column:status" description:"状态"`
}

// TableName ...
func (t *DomainInfo) TableName() string {
	return "Hawkeye_domain"
}

// Find ...
func (t *DomainInfo) Find(c *gin.Context, tx *gorm.DB, search *DomainInfo) (*DomainInfo, error) {
	model := &DomainInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *DomainInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete ...
func (t *DomainInfo) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// AllRecord ...
func (t *DomainInfo) AllRecord(c *gin.Context, tx *gorm.DB) ([]DomainInfo, int64, error) {
	var list []DomainInfo
	var count int64

	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	if t.AssetID > 0 {
		query.Where(fmt.Sprintf(" asset_id=%v", t.AssetID))
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
