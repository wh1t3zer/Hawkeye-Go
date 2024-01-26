package dao

import (
	"fmt"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/dto"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// AssetInfo ...
type AssetInfo struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	TaskID    int64     `json:"task_id" gorm:"column:task_id" description:"任务ID"`
	IP        string    `json:"ip" gorm:"column:ip" description:"IP"`
	GPS       string    `json:"gps" gorm:"column:gps" description:"GPS"`
	AREA      string    `json:"area" gorm:"column:area" description:"区域"`
	ISP       string    `json:"isp" gorm:"column:isp" description:"运营商"`
	OS        string    `json:"os" gorm:"column:os" description:"操作系统"`
	Vendor    string    `json:"vendor" gorm:"column:vendor" description:"设备"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

// TableName ...
func (t *AssetInfo) TableName() string {
	return "Hawkeye_asset"
}

// Find ...
func (t *AssetInfo) Find(c *gin.Context, tx *gorm.DB, search *AssetInfo) (*AssetInfo, error) {
	model := &AssetInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *AssetInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete ...
func (t *AssetInfo) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// PageList ...
func (t *AssetInfo) PageList(c *gin.Context, tx *gorm.DB, params *dto.PublicListInput) ([]AssetInfo, int64, error) {
	var list []AssetInfo
	var count int64
	pageNo := params.Page
	pageSize := params.Limit

	//limit offset,pagesize
	offset := (pageNo - 1) * pageSize
	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	query = query.Where("is_delete=?", 0)
	if params.Info != "" {
		query = query.Where(" (ip like ? or vendor like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
	}
	err := query.Limit(pageSize).Offset(offset).Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// AllRecord 所有记录, 除去被删除之外(接受task_id查询)
func (t *AssetInfo) AllRecord(c *gin.Context, tx *gorm.DB) ([]AssetInfo, int64, error) {
	var list []AssetInfo
	var count int64

	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	query = query.Where("is_delete=?", 0)
	if t.TaskID > 0 {
		query = query.Where(fmt.Sprintf("task_id=%v", t.TaskID))
		err := query.Order("id desc").Find(&list).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}
		if errCount := query.Count(&count).Error; errCount != nil {
			return nil, 0, err
		}
	} else {
		if err := query.Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}
		if errCount := query.Count(&count).Error; errCount != nil {
			return nil, 0, errCount
		}
	}
	return list, count, nil
}
