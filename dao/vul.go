package dao

import (
	"fmt"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/dto"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// VulInfo ...
type VulInfo struct {
	ID              int64     `json:"id" gorm:"primary_key"`
	AssetID         int64     `json:"asset_id" gorm:"column:asset_id" description:"资产id"`
	PortID          int64     `json:"port_id" gorm:"column:port_id" description:"端口id"`
	PluginID        int64     `json:"plugin_id" gorm:"column:plugin_id" description:"插件ID"`
	VerifyURL       string    `json:"verify_url" gorm:"column:verify_url" description:"漏洞验证URL"`
	VerifyPayload   string    `json:"verify_payload" gorm:"column:verify_payload" description:"漏洞验证Payload"`
	VerifyResult    string    `json:"verify_result" gorm:"column:verify_result" description:"漏洞验证Result"`
	ExploitURL      string    `json:"exploit_url" gorm:"column:exploit_url" description:"漏洞利用URL"`
	ExploitPayload  string    `json:"exploit_payload" gorm:"column:exploit_payload" description:"漏洞利用Payload"`
	ExploitResult   string    `json:"exploit_result" gorm:"column:exploit_result" description:"漏洞利用Result"`
	WebshellURL     string    `json:"webshell_url" gorm:"column:webshell_url" description:"Webshell URL"`
	WebshellPayload string    `json:"webshell_payload" gorm:"column:webshell_payload" description:"Webshell Payload"`
	WebshellResult  string    `json:"webshell_result" gorm:"column:webshell_result" description:"Webshell Result"`
	TrojanURL       string    `json:"trojan_url" gorm:"column:trojan_url" description:"Trojan URL"`
	TrojanPayload   string    `json:"trojan_payload" gorm:"column:trojan_payload" description:"Trojan Payload"`
	TrojanResult    string    `json:"trojan_result" gorm:"column:trojan_result" description:"Trojan Result"`
	CreatedAt       time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间"`
	IsDelete        int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

// TableName ...
func (t *VulInfo) TableName() string {
	return "Hawkeye_vulinfo"
}

// Find ...
func (t *VulInfo) Find(c *gin.Context, tx *gorm.DB, search *VulInfo) (*VulInfo, error) {
	model := &VulInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *VulInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete ...
func (t *VulInfo) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// PageList ...
func (t *VulInfo) PageList(c *gin.Context, tx *gorm.DB, params *dto.PublicListInput) ([]VulInfo, int64, error) {
	var list []VulInfo
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
func (t *VulInfo) AllRecord(c *gin.Context, tx *gorm.DB) ([]VulInfo, int64, error) {
	var list []VulInfo
	var count int64

	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	if t.PortID > 0 {
		query = query.Where(fmt.Sprintf("port_id=%v", t.PortID))
		err := query.Order("id desc").Find(&list).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}
		errCount := query.Count(&count).Error
		if errCount != nil {
			return nil, 0, err
		}
	} else {
		err := query.Order("id desc").Find(&list).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}
		errCount := query.Count(&count).Error
		if errCount != nil {
			return nil, 0, err
		}
	}
	return list, count, nil
}
