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
	Id                  string               `json:"id"`
	HostName            string               `json:"hostname"`
	Protocol            string               `json:"protocol"`
	Servers             []Server             `json:"server"`
	Proxy               bool                 `json:"proxy"`
	Locked              int                  `json:"locked"`
	Timestamp           int64                `json:"timestamp"`
	Tls                 string               `json:"tls"`
	Cipher              string               `json:"cipher"`
	Extend              map[string]string    `json:"extend"`
	Flag                map[string]string    `json:"flag"`
	Description         string               `json:"description"`
	PolicyId            string               `json:"policyid"`
	DomainId            string               `json:"domainid"`
	ProjectId           string               `json:"projectid"`
	EnterpriseProjectId string               `json:"enterprise_project_id"`
	CertificateId       string               `json:"certificateid"`
	CertificateName     string               `json:"certificatename"`
	ProtectStatus       int                  `json:"protect_status"`
	AccessStatus        int                  `json:"access_status"`
	WebTag              string               `json:"web_tag"`
	BlockPage           DomainBlockPage      `json:"block_page"`
	TrafficMark         DomainTrafficMark    `json:"traffic_mark"`
	CircuitBreaker      DomainCircuitBreaker `json:"circuit_breaker"`
	TimeoutConfig       DomainTimeoutConfig  `json:"timeout_config"`
	ForwardHeaderMap    map[string]string    `json:"forward_header_map"`
	AccessProgress      []AccessProgress     `json:"access_progress"`

	// Deprecated
	AccessCode string `json:"access_code"`
	// Deprecated
	Mode string `json:"mode"`
	// Deprecated
	PoolIds []string `json:"pool_ids"`
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

type DomainCircuitBreaker struct {
	Switch           bool    `json:"switch"`
	DeadNum          int     `json:"dead_num"`
	DeadRatio        float64 `json:"dead_ratio"`
	BlockTime        int     `json:"block_time"`
	SuperpositionNum int     `json:"superposition_num"`
	SuspendNum       int     `json:"suspend_num"`
	SusBlockTime     int     `json:"sus_block_time"`
}

type DomainTimeoutConfig struct {
	ConnectTimeout int `json:"connect_timeout"`
	SendTimeout    int `json:"send_timeout"`
	ReadTimeout    int `json:"read_timeout"`
}

type AccessProgress struct {
	Step   int `json:"step"`
	Status int `json:"status"`
}

type PremiumHostList struct {
	Total int                 `json:"total"`
	Items []SimplePremiumHost `json:"items"`
}

type PremiumHostProtectStatus struct {
	KeepPolicy bool
}
