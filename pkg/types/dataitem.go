package types

import "net/http"

type DataItem struct {
	Id             string
	Domain         string
	Url            string
	Method         string
	Https          bool //http/https flag
	SourceRequest  *http.Request
	SourceResponse *http.Response
	VulnType       []string
	VulnRequest    []*http.Request
	VulnResponse   []*http.Response
	ReportTime     string
	CheckState     bool
}
