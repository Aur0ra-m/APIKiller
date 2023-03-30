package types

type DataItemStr struct {
	Id             string `json:"Id" form:"Id" `
	Domain         string `json:"Domain" form:"Domain" `
	Url            string `json:"Url" form:"Url" `
	Method         string `json:"Method" form:"Method"`
	Https          bool   `json:"Https" form:"Https" `
	SourceRequest  string `json:"SourceRequest" form:"SourceRequest" `
	SourceResponse string `json:"SourceResponse" form:"SourceResponse" `
	VulnType       string `json:"VulnType" form:"VulnType" `
	VulnRequest    string `json:"VulnRequest" form:"VulnRequest" `
	VulnResponse   string `json:"VulnResponse" form:"VulnResponse" `
	ReportTime     string `json:"ReportTime" form:"ReportTime" `
	CheckState     bool   `json:"CheckState" form:"CheckState" `
}
