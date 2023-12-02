package dto

// WebInfoOutput web信息
type WebInfoOutput struct {
	ID           int64  `json:"id" form:"primary_key"`
	PortID       int64  `json:"port_id" form:"port_id" description:"端口id"`
	StartURL     string `json:"start_url" form:"start_url" description:"起始URL"`
	Title        string `json:"title" form:"title" description:"站点标题"`
	Server       string `json:"server" form:"server" description:"Web服务器"`
	ContentType  string `json:"content_type" form:"content_type" description:"内容类型"`
	LoginList    string `json:"login_list" form:"login_list" description:"登录页列表"`
	UploadList   string `json:"upload_list" form:"upload_list" description:"上传页面列表"`
	SubDomain    string `json:"sub_domain" form:"sub_domain" description:"子域名列表"`
	RouteList    string `json:"route_list" form:"route_list" description:"URL列表"`
	ResourceList string `json:"resource_list" form:"resource_list" description:"资源列表"`
}

// WebListOutput ...
type WebListOutput struct {
	Total int64            `json:"total" form:"total" comment:"总数" example:"" validate:""` //总数
	List  []*WebInfoOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   //列表
}
