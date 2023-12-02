package dao

import (
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye/dto"
	"github.com/wh1t3zer/Hawkeye/utils"
)

// PocTask ...
type PocTask struct {
	ID         int       `json:"id" gorm:"primary_key" description:"自增主键"`
	AssetID    int       `json:"asset_id" gorm:"column:asset_id" description:"资产ID"`
	PortinfoID int       `json:"portinfo_id" gorm:"column:portinfo_id" description:"端口ID"`
	Status     string    `json:"status" gorm:"column:status" description:"任务状态"`
	TargetList string    `json:"target_list" gorm:"column:target_list" description:"目标列表"`
	TaskName   string    `json:"task_name" gorm:"column:task_name" description:"任务名"`
	Recursion  uint8     `json:"recursion" gorm:"column:recursion" description:"扫描周期"`
	PluginList string    `json:"plugin_list" gorm:"column:plugin_list" description:"插件列表"`
	UpdatedAt  time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt  time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete   int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

// TableName ...
func (t *PocTask) TableName() string {
	return "poc_task"
}

// Find ...
func (t *PocTask) Find(c *gin.Context, tx *gorm.DB, search *PocTask) (*PocTask, error) {
	out := &PocTask{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Save ...
func (t *PocTask) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error
}

// PageList 分页查询=>得服务数组
func (t *PocTask) PageList(c *gin.Context, tx *gorm.DB, search *dto.PublicListInput) ([]PocTask, int64, error) {
	var total int64 = 0
	list := []PocTask{}
	offset := (search.Page - 1) * search.Limit //第一页的话就不用偏移了,直接limit查
	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Where("is_delete=0")
	if search.Info != "" { ///模糊查询
		query = query.Where("name like ?", "%"+search.Info+"%")
	}
	if err := query.Limit(search.Limit).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(search.Page).Offset(offset).Count(&total)
	return list, total, nil
}
