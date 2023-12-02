package dto

// PanelGroupData 头顶四个数据统计box
type PanelGroupData struct {
	VulCount      int `json:"vul_count"`
	AssetCount    int `json:"asset_count"`
	ServiceCount  int `json:"service_count"`
	ResourceCount int `json:"resource_count"`
}

// ChartSeries chart数据源
type ChartSeries struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// TableSeries 表格数据源
type TableSeries struct {
	VulID   string `json:"vul_id"`   // 漏洞ID
	VulName string `json:"vul_name"` // 漏洞名
	VulType string `json:"vul_type"` // 漏洞类型
}

// TimeLineSeries 时间线数据源
type TimeLineSeries struct {
	Content   string `json:"content"`   // 进度内容 漏洞渗透
	Timestamp string `json:"timestamp"` // 时间非时间戳 2018-04-12 20:46
	Icon      string `json:"icon"`      // el-icon-check、el-icon-loading
	Type      string `json:"type"`      // primary(蓝色) success(绿色)
}

// ChartBoxCard 卡片内容
type ChartBoxCard struct {
	Title  string        `json:"title"`  // 卡片标题
	Image  string        `json:"image"`  // 卡片头像
	Type   string        `json:"type"`   // 卡片绘图类型 pie line tatle
	Series []ChartSeries `json:"series"` // 数据源
}

// TableBoxCard 卡片内容
type TableBoxCard struct {
	Title  string          `json:"title"`  // 卡片标题
	Image  string          `json:"image"`  // 卡片头像
	Type   string          `json:"type"`   // 卡片绘图类型 pie line tatle
	Series []VulInfoOutput `json:"series"` // 数据源
}

// TimeLineBoxCard 卡片内容
type TimeLineBoxCard struct {
	Title  string           `json:"title"`  // 卡片标题
	Image  string           `json:"image"`  // 卡片头像
	Type   string           `json:"type"`   // 卡片绘图类型 pie line tatle
	Series []TimeLineSeries `json:"series"` // 数据源
}

// Vulnerability 漏洞信息
type Vulnerability struct {
	VulID   string `json:"vul_id"`
	VulName string `json:"vul_name"`
	VulType string `json:"vul_type"`
}

// DashboardOutput 全局
type DashboardOutput struct {
	PanelGroup *PanelGroupData `json:"panel_data"`
	Box1       *ChartBoxCard   `json:"box1"`
	Box2       *ChartBoxCard   `json:"box2"`
	Box3       *ChartBoxCard   `json:"box3"`
	Box4       *ChartBoxCard   `json:"box4"`
	Box5       *TableBoxCard   `json:"box5"`
	Box6       *ChartBoxCard   `json:"box6"`
}

// TaskDashboardOutput 任务视图
type TaskDashboardOutput struct {
	PanelGroup *PanelGroupData `json:"panel_data"`
	Box1       *ChartBoxCard   `json:"box1"`
	Box2       *ChartBoxCard   `json:"box2"`
	Box3       *ChartBoxCard   `json:"box3"`
	Box4       *ChartBoxCard   `json:"box4"`
	Box5       *TableBoxCard   `json:"box5"`
	Box6       *ChartBoxCard   `json:"box6"`
	Status     string          `json:"status"`
	Percent    int8            `json:"percent"`
}
