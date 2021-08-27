/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package premium_domains

type CreatePremiumHostRst struct {
	Id        string `json:"id"`
	PolicyId  string `json:"policyid"`
	HostName  string `json:"hostname"`
	DomainId  string `json:"domainid"`
	ProjectId string `json:"projectid"`
	Protocol  string `json:"protocol"`
}

type PremiumHost struct {
	Id              string            `json:"id"`
	PolicyId        string            `json:"policyid"`
	HostName        string            `json:"hostname"`
	DomainId        string            `json:"domainid"`
	ProjectId       string            `json:"project_id"`
	AccessCode      string            `json:"access_code"`
	Protocol        string            `json:"protocol"`
	Servers         []Server          `json:"server"`
	CertificateId   string            `json:"certificateid"`
	CertificateName string            `json:"certificatename"`
	Tls             string            `json:"tls"`
	Cipher          string            `json:"cipher"`
	Proxy           bool              `json:"proxy"`
	Locked          int               `json:"locked"`
	ProtectStatus   int               `json:"protect_status"`
	AccessStatus    int               `json:"access_status"`
	Timestamp       int64             `json:"timestamp"`
	BlockPage       DomainBlockPage   `json:"block_page"`
	Extend          map[string]string `json:"extend"`
	TrafficMark     DomainTrafficMark `json:"traffic_mark"`
	Flag            map[string]string `json:"flag"`
	Mode            string            `json:"mode"`
	PoolIds         []string          `json:"pool_ids"`
}

type SimplePremiumHost struct {
	Id            string            `json:"id"`
	Hostname      string            `json:"hostname"`
	PolicyId      string            `json:"policyid"`
	ProtectStatus int               `json:"protect_status"`
	AccessStatus  int               `json:"access_status"`
	Flag          map[string]string `json:"flag"`
	Mode          string            `json:"mode"`
	PoolIds       []string          `json:"pool_ids"`
}

type DomainBlockPage struct {
	Template    string           `json:"template"`
	CustomPage  DomainCustomPage `json:"custom_page"`
	RedirectUrl string           `json:"redirect_url"`
}

type DomainCustomPage struct {
	StatusCode  string `json:"status_code"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
}

type DomainTrafficMark struct {
	Sip    []string `json:"sip"`
	Cookie string   `json:"cookie"`
	Params string   `json:"params"`
}

type PremiumHostList struct {
	Total int                 `json:"total"`
	Items []SimplePremiumHost `json:"items"`
}

type PremiumHostProtectStatus struct {
	KeepPolicy bool
}
