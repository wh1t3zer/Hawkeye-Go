package dao

import (
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/dto"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// TrapInfo ...
type TrapInfo struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	AssetID   int64     `json:"asset_id" gorm:"column:asset_id" description:"资产id"`
	PortID    int64     `json:"port_id" gorm:"column:port_id" description:"端口id"`
	PluginID  int64     `json:"plugin_id" gorm:"column:plugin_id" description:"插件ID"`
	Verify    string    `json:"verify" gorm:"column:verify" description:"验证项"`
	TrapID    string    `json:"trap_id" gorm:"column:trap_id" description:"蜜罐ID"`
	Name      string    `json:"name" gorm:"column:name" description:"插件名"`
	Protocol  string    `json:"protocol" gorm:"column:protocol" description:"协议"`
	AppName   string    `json:"app_name" gorm:"column:app_name" description:"应用名"`
	HoneyPot  string    `json:"honeypot" gorm:"column:honeypot" description:"蜜罐名"`
	Desc      string    `json:"desc" gorm:"column:desc" description:"描述"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间"`
}

// TableName ...
func (t *TrapInfo) TableName() string {
	return "Hawkeye_trap"
}

// Find ...
func (t *TrapInfo) Find(c *gin.Context, tx *gorm.DB, search *TrapInfo) (*TrapInfo, error) {
	model := &TrapInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *TrapInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete ...
func (t *TrapInfo) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// PageList ...
func (t *TrapInfo) PageList(c *gin.Context, tx *gorm.DB, params *dto.PublicListInput) ([]TrapInfo, int64, error) {
	var list []TrapInfo
	var count int64
	pageNo := params.Page
	pageSize := params.Limit

	offset := (pageNo - 1) * pageSize
	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	if params.Info != "" {
		query = query.Where(" (id like ? or vendor like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
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
func (t *TrapInfo) AllRecord(condition string, c *gin.Context, tx *gorm.DB) ([]TrapInfo, int64, error) {
	// condition为where的查询命令 port_id=1 AND asset_id=2
	var list []TrapInfo
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
