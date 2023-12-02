package dao

import (
	"fmt"
	"strings"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye/dto"
	"github.com/wh1t3zer/Hawkeye/utils"
)

// TaskInfo ...
type TaskInfo struct {
	ID         int64     `json:"id" gorm:"primary_key"`
	RuleID     int64     `json:"rule_id" gorm:"column:rule_id" description:"web规则id"`
	Name       string    `json:"name" gorm:"column:name" description:"任务名"`
	TargetList string    `json:"target_list" gorm:"column:target_list" description:"目标列表"`
	WebScan    int8      `json:"web_scan" gorm:"column:web_scan" description:"Web扫描"`
	PocScan    int8      `json:"poc_scan" gorm:"column:poc_scan" description:"Poc扫描"`
	AuthScan   int8      `json:"auth_scan" gorm:"column:auth_scan" description:"权限扫描"`
	TrapScan   int8      `json:"trap_scan" gorm:"column:trap_scan" description:"蜜罐识别"`
	Recursion  int8      `json:"recursion" gorm:"column:recursion" description:"扫描周期"`
	Progress   string    `json:"progress" gorm:"column:progress" description:"扫描进程"`
	Percent    int8      `json:"percent" gorm:"column:percent" description:"扫描百分比0-100"`
	Status     string    `json:"status" gorm:"column:status" description:"扫描状态"`
	CreatedAt  time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间"`
	UpdatedAt  time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete   int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

// TableName ...
func (t *TaskInfo) TableName() string {
	return "Hawkeye_task"
}

// Find ...
func (t *TaskInfo) Find(c *gin.Context, tx *gorm.DB, search *TaskInfo) (*TaskInfo, error) {
	model := &TaskInfo{}
	err := tx.SetCtx(utils.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *TaskInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(utils.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// PageList 模糊分页查询
func (t *TaskInfo) PageList(c *gin.Context, tx *gorm.DB, params *dto.PublicListInput) ([]TaskInfo, int64, error) {
	var list []TaskInfo
	var count int64
	pageNo := params.Page
	pageSize := params.Limit

	//limit offset,pagesize
	offset := (pageNo - 1) * pageSize
	query := tx.SetCtx(utils.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	query = query.Where("is_delete=?", 0)
	if params.Info != "" {
		query = query.Where(" (name like ? or target_list like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
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

// TaskStatistics 根据任务ID 找 资产详情
func (t *TaskInfo) TaskStatistics(c *gin.Context) *dto.TaskDashboardOutput {
	// 0、 查看当前的任务状态
	taskinfo, err := t.Find(c, lib.GORMDefaultPool, t)
	if err != nil {
		fmt.Println("null--->", taskinfo)
		return nil
	}
	out := &dto.TaskDashboardOutput{
		PanelGroup: &dto.PanelGroupData{},
		Box1: &dto.ChartBoxCard{
			Title: "Hardware", Image: "https://zan71.com/cdn-img/icon/avatar/tx1.gif", Type: "pie", Series: []dto.ChartSeries{},
		},
		Box2: &dto.ChartBoxCard{
			Title: "Software", Image: "https://zan71.com/cdn-img/icon/avatar/tx.gif", Type: "pie", Series: []dto.ChartSeries{},
		},
		Box3: &dto.ChartBoxCard{
			Title: "SubDomain", Image: "https://zan71.com/cdn-img/icon/avatar/tx.gif", Type: "pie", Series: []dto.ChartSeries{},
		},
		Box4: &dto.ChartBoxCard{
			Title: "Vul Type", Image: "https://zan71.com/cdn-img/icon/avatar/tx1.gif", Type: "pie", Series: []dto.ChartSeries{},
		},
		Box5: &dto.TableBoxCard{
			Title: "Latest Vul", Image: "https://zan71.com/cdn-img/icon/avatar/tx1.gif", Type: "table", Series: []dto.VulInfoOutput{},
		},
		Box6: &dto.ChartBoxCard{
			Title: "Web Site", Image: "https://zan71.com/cdn-img/icon/avatar/tx.gif", Type: "line", Series: []dto.ChartSeries{},
		},
		Status: taskinfo.Status, Percent: taskinfo.Percent,
	}

	// 1、 根据任务ID找旗下所有资产
	search := &AssetInfo{TaskID: t.ID}
	assetArray, assetCount, err := search.AllRecord(c, lib.GORMDefaultPool)
	if err != nil {
		return out
	}

	var srvCount = 0
	var vulCount = 0
	var vendors = make(map[string]int, 10)
	var domains = make(map[string]int, 10)
	var softwares = make(map[string]int, 10)
	var webinfos = make(map[string]int, 10)
	var vultypes = make(map[string]int, 10)
	var outarray = []dto.VulInfoOutput{}

	webinfos["Web Server"] = 0
	webinfos["Content Type"] = 0
	webinfos["Login Page"] = 0
	webinfos["Upload Page"] = 0
	webinfos["Sub Domain"] = 0
	webinfos["Site URL"] = 0
	webinfos["Resource Path"] = 0
	// 2、 遍历所有资产
	for _, asset := range assetArray {
		vendors[asset.Vendor]++
		// 2.1 根据资产ID找域名表
		search1 := &DomainInfo{AssetID: asset.ID}
		search1, _ = search1.Find(c, lib.GORMDefaultPool, search1)

		for _, item := range strings.Split(search1.SubDomainList, ",") {
			if item != "" {
				domains[item]++
			}
		}

		// 2.2 根据资产ID找端口表
		search2 := &PortInfo{AssetID: asset.ID}
		srvArray, total, _ := search2.AllRecord(c, lib.GORMDefaultPool)
		srvCount += int(total)
		// 2.2.1 遍历端口表
		for _, portinfo := range srvArray {
			software := portinfo.Product
			if software == "" {
				software = portinfo.Name
			}
			softwares[software]++
			// 2.2.1.1 根据端口ID找Web信息
			web := &WebInfo{PortID: portinfo.ID}
			if webArray, _, err := web.AllRecord(c, lib.GORMDefaultPool); err == nil {
				for _, info := range webArray {
					if strings.TrimSpace(info.Server) != "" {
						webinfos["Web Server"] += len(strings.Split(info.Server, ","))
					}
					if strings.TrimSpace(info.ContentType) != "" {
						webinfos["Content Type"] += len(strings.Split(info.ContentType, ","))
					}
					if strings.TrimSpace(info.LoginList) != "" {
						webinfos["Login Page"] += len(strings.Split(info.LoginList, ","))
					}
					if strings.TrimSpace(info.UploadList) != "" {
						webinfos["Upload Page"] += len(strings.Split(info.UploadList, ","))
					}
					if strings.TrimSpace(info.SubDomain) != "" {
						webinfos["Sub Domain"] += len(strings.Split(info.SubDomain, ","))
					}
					if strings.TrimSpace(info.RouteList) != "" {
						webinfos["Site URL"] += len(strings.Split(info.RouteList, ","))
					}
					if strings.TrimSpace(info.ResourceList) != "" {
						webinfos["Resource Path"] += len(strings.Split(info.ResourceList, ","))
					}
				}
			}
			// 2.2.1.2 根据端口ID找漏洞信息
			vul := &VulInfo{PortID: portinfo.ID}
			if vulArray, vtotal, err := vul.AllRecord(c, lib.GORMDefaultPool); err == nil {
				vulCount += int(vtotal)
				for _, info := range vulArray {
					// 漏洞类型
					poctObj := &PocPlugin{ID: info.PluginID}
					poc, _ := poctObj.Find(c, lib.GORMDefaultPool, poctObj)
					vultypes[poc.VulType]++
					// 漏洞列表
					outarray = append(outarray, dto.VulInfoOutput{
						ID:              info.ID,
						AssetID:         info.AssetID,
						Asset:           fmt.Sprintf("%v:%v", asset.IP, portinfo.Port),
						PortID:          info.PortID,
						PluginID:        info.PluginID,
						AppName:         poc.AppName,
						VulName:         poc.VulName,
						VulType:         poc.VulType,
						VerifyURL:       info.VerifyURL,
						VerifyPayload:   info.VerifyPayload,
						VerifyResult:    info.VerifyResult,
						ExploitURL:      info.ExploitURL,
						ExploitPayload:  info.ExploitPayload,
						ExploitResult:   info.ExploitResult,
						WebshellURL:     info.WebshellURL,
						WebshellPayload: info.WebshellPayload,
						WebshellResult:  info.WebshellResult,
						TrojanURL:       info.TrojanURL,
						TrojanPayload:   info.TrojanPayload,
						TrojanResult:    info.TrojanResult,
						CreatedAt:       info.CreatedAt.Format("2006-01-02 15:04:05"),
						IsDelete:        info.IsDelete,
					})
				}
			}
		}
	}
	// [*]组合输出数据
	out.PanelGroup.AssetCount = int(assetCount)
	out.PanelGroup.ServiceCount = int(srvCount)
	out.PanelGroup.VulCount = vulCount
	out.Box5.Series = outarray
	for key, value := range vendors {
		out.Box1.Series = append(out.Box1.Series, dto.ChartSeries{Name: key, Value: value})
	}
	for key, value := range softwares {
		out.Box2.Series = append(out.Box2.Series, dto.ChartSeries{Name: key, Value: value})
	}
	for key, value := range domains {
		out.Box3.Series = append(out.Box3.Series, dto.ChartSeries{Name: key, Value: value})
	}
	for key, value := range vultypes {
		out.Box4.Series = append(out.Box4.Series, dto.ChartSeries{Name: key, Value: value})
	}
	for key, value := range webinfos {
		out.Box6.Series = append(out.Box6.Series, dto.ChartSeries{Name: key, Value: value})
	}
	return out
}
