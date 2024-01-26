package dao

import (
	"fmt"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/dto"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// WebInfo ...
type WebInfo struct {
	ID           int64  `json:"id" gorm:"primary_key"`
	PortID       int64  `json:"port_id" gorm:"column:port_id" description:"端口id"`
	StartURL     string `json:"start_url" gorm:"column:start_url" description:"起始URL"`
	Title        string `json:"title" gorm:"column:title" description:"站点标题"`
	Server       string `json:"server" gorm:"column:server" description:"Web服务器"`
	ContentType  string `json:"content_type" gorm:"column:content_type" description:"内容类型"`
	LoginList    string `json:"login_list" gorm:"column:login_list" description:"登录页列表"`
	UploadList   string `json:"upload_list" gorm:"column:upload_list" description:"上传页面列表"`
	SubDomain    string `json:"sub_domain" gorm:"column:sub_domain" description:"子域名列表"`
	RouteList    string `json:"route_list" gorm:"column:route_list" description:"URL列表"`
	ResourceList string `json:"resource_list" gorm:"column:resource_list" description:"资源列表"`
}

// TableName ...
func (t *WebInfo) TableName() string {
	return "Hawkeye_webinfo"
}

// Find ...
func (t *WebInfo) Find(c *gin.Context, tx *gorm.DB, search *WebInfo) (*WebInfo, error) {
	model := &WebInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *WebInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete ...
func (t *WebInfo) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// PageList ...
func (t *WebInfo) PageList(c *gin.Context, tx *gorm.DB, params *dto.PublicListInput) ([]WebInfo, int64, error) {
	var list []WebInfo
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
func (t *WebInfo) AllRecord(c *gin.Context, tx *gorm.DB) ([]WebInfo, int64, error) {
	var list []WebInfo
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
