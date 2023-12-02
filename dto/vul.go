package dto

// VulInfoOutput 漏洞信息
type VulInfoOutput struct {
	ID              int64  `json:"id" gorm:"primary_key"`
	AssetID         int64  `json:"asset_id" form:"asset_id" description:"资产id"`
	Asset           string `json:"asset" form:"asset" description:"资产目标"`
	PortID          int64  `json:"port_id" form:"port_id" description:"端口id"`
	PluginID        int64  `json:"plugin_id" form:"plugin_id" description:"插件ID"`
	AppName         string `json:"app_name" form:"app_name" description:"应用名"`
	VulName         string `json:"vul_name" form:"vul_name" description:"漏洞名"`
	VulType         string `json:"vul_type" form:"vul_type" description:"漏洞类型"`
	VerifyURL       string `json:"verify_url" form:"verify_url" description:"漏洞验证URL"`
	VerifyPayload   string `json:"verify_payload" form:"verify_payload" description:"漏洞验证Payload"`
	VerifyResult    string `json:"verify_result" form:"verify_result" description:"漏洞验证Result"`
	ExploitURL      string `json:"exploit_url" form:"exploit_url" description:"漏洞利用URL"`
	ExploitPayload  string `json:"exploit_payload" form:"exploit_payload" description:"漏洞利用Payload"`
	ExploitResult   string `json:"exploit_result" form:"exploit_result" description:"漏洞利用Result"`
	WebshellURL     string `json:"webshell_url" form:"webshell_url" description:"Webshell URL"`
	WebshellPayload string `json:"webshell_payload" form:"webshell_payload" description:"Webshell Payload"`
	WebshellResult  string `json:"webshell_result" form:"webshell_result" description:"Webshell Result"`
	TrojanURL       string `json:"trojan_url" form:"trojan_url" description:"Trojan URL"`
	TrojanPayload   string `json:"trojan_payload" form:"trojan_payload" description:"Trojan Payload"`
	TrojanResult    string `json:"trojan_result" form:"trojan_result" description:"Trojan Result"`
	CreatedAt       string `json:"create_at" form:"create_at" description:"添加时间"`
	IsDelete        int8   `json:"is_delete" form:"is_delete" description:"是否已删除；0：否；1：是"`
	SpareLine       int8   `json:"line" form:"line" description:"是否已删除；0：其他漏洞；1：穿透; 2: 直连"`
}

// VulListOutput ...
type VulListOutput struct {
	Total int64            `json:"total" form:"total" comment:"总数" example:"" validate:""` //总数
	List  []*VulInfoOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   //列表
}
