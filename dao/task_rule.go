package dao

import (
	"time"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/dto"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// RuleInfo ...
// type RuleInfo struct {
// 	ID        int64  `json:"id" gorm:"primary_key"`
// 	AssetID   int64  `json:"asset_id" gorm:"column:asset_id" description:"资产ID"`
// 	Port      string `json:"port" gorm:"column:port" description:"端口"`
// 	Name      string `json:"name" gorm:"column:name" description:"服务名"`
// 	State     string `json:"state" gorm:"column:state" description:"状态"`
// 	Product   string `json:"product" gorm:"column:product" description:"应用"`
// 	Version   string `json:"version" gorm:"column:version" description:"版本"`
// 	Extrainfo string `json:"extrainfo" gorm:"column:extrainfo" description:"服务名"`
// 	Conf      string `json:"conf" gorm:"column:conf" description:"未知"`
// 	Cpe       string `json:"cpe" gorm:"column:cpe" description:"指纹"`
// }

//ruleinfo ...
type RuleInfo struct{
	ID             int64  `json:"id" gorm:"primary_key"`
	Trojan_cmd     string  `json:"trojan_cmd" gorm:"column:asset_id" description:"木马"`
	Port_list      string `json:"port_list" gorm:"column:port_list" description:"端口字典"`
	Domain_dict    string `json:"domain_dict" gorm:"column:domain_dict" description:"域名字典"`
	User_dict      string `json:"user_dict" gorm:"column:user_dict" description:"用户名字典"`
	Passwd_dict    string `json:"passwd_dict" gorm:"column:passwd_dict" description:"密码字典"`
	Create_at      time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`

}

// TableName ...
func (t *RuleInfo) TableName() string {
	return "task_rule"
}

// Find ...
func (t *RuleInfo) Find(c *gin.Context, tx *gorm.DB, search *RuleInfo) (*RuleInfo, error) {
	model := &RuleInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *RuleInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete ...
func (t *RuleInfo) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// PageList ...
func (t *RuleInfo) PageList(c *gin.Context, tx *gorm.DB, params *dto.PublicListInput) ([]RuleInfo, int64, error) {
	var list []RuleInfo
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
func (t *RuleInfo) AllRecord(condition string, c *gin.Context, tx *gorm.DB) ([]RuleInfo, int64, error) {
	var list []RuleInfo
	var count int64

	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	if condition != "" {
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
