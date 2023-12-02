package dao

import (
	"fmt"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye/dto"
	"github.com/wh1t3zer/Hawkeye/utils"
)

// PortInfo ...
type PortInfo struct {
	ID        int64  `json:"id" gorm:"primary_key"`
	AssetID   int64  `json:"asset_id" gorm:"column:asset_id" description:"资产ID"`
	Port      string `json:"port" gorm:"column:port" description:"端口"`
	Name      string `json:"name" gorm:"column:name" description:"服务名"`
	State     string `json:"state" gorm:"column:state" description:"状态"`
	Product   string `json:"product" gorm:"column:product" description:"应用"`
	Version   string `json:"version" gorm:"column:version" description:"版本"`
	Extrainfo string `json:"extrainfo" gorm:"column:extrainfo" description:"服务名"`
	Conf      string `json:"conf" gorm:"column:conf" description:"未知"`
	Cpe       string `json:"cpe" gorm:"column:cpe" description:"指纹"`
}

// TableName ...
func (t *PortInfo) TableName() string {
	return "Hawkeye_portinfo"
}

// Find ...
func (t *PortInfo) Find(c *gin.Context, tx *gorm.DB, search *PortInfo) (*PortInfo, error) {
	model := &PortInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *PortInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete ...
func (t *PortInfo) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// PageList ...
func (t *PortInfo) PageList(c *gin.Context, tx *gorm.DB, params *dto.PublicListInput) ([]PortInfo, int64, error) {
	var list []PortInfo
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
func (t *PortInfo) AllRecord(c *gin.Context, tx *gorm.DB) ([]PortInfo, int64, error) {
	var list []PortInfo
	var count int64

	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	if t.AssetID > 0 {
		query = query.Where(fmt.Sprintf("asset_id=%v", t.AssetID))
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
