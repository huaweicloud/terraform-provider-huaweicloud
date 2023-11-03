/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package premium_domains

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/utils"
)

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreatePremiumHostOpts the options for creating premium domains.
type CreateOpts struct {
	CertificateId       string            `json:"certificateid,omitempty"`
	CertificateName     string            `json:"certificatename,omitempty"`
	HostName            string            `json:"hostname" required:"true"`
	Proxy               *bool             `json:"proxy,omitempty"`
	PolicyId            string            `json:"policyid,omitempty"`
	Servers             []Server          `json:"server,omitempty"`
	BlockPage           *BlockPage        `json:"block_page,omitempty"`
	ForwardHeaderMap    map[string]string `json:"forward_header_map,omitempty"`
	Description         string            `json:"description,omitempty"`
	EnterpriseProjectID string            `q:"enterprise_project_id" json:"-"`
}

type BlockPage struct {
	Template    string      `json:"template,omitempty"`
	CustomPage  *CustomPage `json:"custom_page,omitempty"`
	RedirectUrl string      `json:"redirect_url,omitempty"`
}

type CustomPage struct {
	StatusCode  string `json:"status_code,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Content     string `json:"content,omitempty"`
}

// PremiumDomainServer the options of domain server for creating premium domains.
type Server struct {
	FrontProtocol string `json:"front_protocol" required:"true"`
	BackProtocol  string `json:"back_protocol" required:"true"`
	Address       string `json:"address" required:"true"`
	Port          int    `json:"port" required:"true"`
	Type          string `json:"type,omitempty"`
	VpcId         string `json:"vpc_id,omitempty"`
}

// Create create a premium domain in HuaweiCloud.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*CreatePremiumHostRst, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c)+query.String(), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r CreatePremiumHostRst
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// Get get a premium domain by id.
func Get(c *golangsdk.ServiceClient, hostId string) (*PremiumHost, error) {
	return GetWithEpsID(c, hostId, "")
}

func GetWithEpsID(c *golangsdk.ServiceClient, hostId, epsID string) (*PremiumHost, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, hostId)+utils.GenerateEpsIDQuery(epsID), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r PremiumHost
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// ListPremiumHostOpts the options for querying a list of premium domains.
type ListPremiumHostOpts struct {
	Page                string `q:"page"`
	PageSize            string `q:"pagesize"`
	HostName            string `q:"hostname"`
	PolicyName          string `q:"policyname"`
	ProtectStatus       int    `q:"protect_status"`
	EnterpriseProjectID string `q:"enterprise_project_id"`
}

// List query a list of premium domains.
func List(c *golangsdk.ServiceClient, opts ListPremiumHostOpts) (*PremiumHostList, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.Get(url, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r PremiumHostList
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// UpdatePremiumHostOpts the options for updating premium domains.
type UpdatePremiumHostOpts struct {
	Proxy               *bool             `json:"proxy,omitempty"`
	CertificateId       string            `json:"certificateid,omitempty"`
	CertificateName     string            `json:"certificatename,omitempty"`
	Tls                 string            `json:"tls,omitempty"`
	Cipher              string            `json:"cipher,omitempty"`
	Description         *string           `json:"description,omitempty"`
	WebTag              *string           `json:"web_tag,omitempty"`
	BlockPage           *BlockPage        `json:"block_page,omitempty"`
	TrafficMark         *TrafficMark      `json:"traffic_mark,omitempty"`
	CircuitBreaker      *CircuitBreaker   `json:"circuit_breaker,omitempty"`
	TimeoutConfig       *TimeoutConfig    `json:"timeout_config,omitempty"`
	Flag                *Flag             `json:"flag,omitempty"`
	ForwardHeaderMap    map[string]string `json:"forward_header_map,omitempty"`
	EnterpriseProjectID string            `q:"enterprise_project_id" json:"-"`
}

type TrafficMark struct {
	Sip    []string `json:"sip,omitempty"`
	Cookie string   `json:"cookie,omitempty"`
	Params string   `json:"params,omitempty"`
}

type CircuitBreaker struct {
	Switch           *bool    `json:"switch,omitempty"`
	DeadNum          *int     `json:"dead_num,omitempty"`
	DeadRatio        *float64 `json:"dead_ratio,omitempty"`
	BlockTime        *int     `json:"block_time,omitempty"`
	SuperpositionNum *int     `json:"superposition_num,omitempty"`
	SuspendNum       *int     `json:"suspend_num,omitempty"`
	SusBlockTime     *int     `json:"sus_block_time,omitempty"`
}

type TimeoutConfig struct {
	ConnectTimeout *int `json:"connect_timeout,omitempty"`
	SendTimeout    *int `json:"send_timeout,omitempty"`
	ReadTimeout    *int `json:"read_timeout,omitempty"`
}

type Flag struct {
	Pci3ds string `json:"pci_3ds,omitempty"`
	PciDss string `json:"pci_dss,omitempty"`
}

// Update update premium domains according to UpdatePremiumHostOpts.
func Update(c *golangsdk.ServiceClient, hostId string, opts UpdatePremiumHostOpts) (*PremiumHost, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, hostId)+query.String(), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r PremiumHost
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// updateProtectStatusOpts the struct for updating the protect status of premium domain.
// Only used in the package.
type updateProtectStatusOpts struct {
	ProtectStatus *int `json:"protect_status" required:"true"`
}

// UpdateProtectStatus update the protect status of premium domain.
func UpdateProtectStatus(c *golangsdk.ServiceClient, hostId string, protectStatus int) (*PremiumHostProtectStatus, error) {
	return UpdateProtectStatusWithWpsID(c, protectStatus, hostId, "")
}

func UpdateProtectStatusWithWpsID(c *golangsdk.ServiceClient, protectStatus int,
	hostId, epsID string) (*PremiumHostProtectStatus, error) {
	opts := updateProtectStatusOpts{
		ProtectStatus: &protectStatus,
	}

	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(protectStatusURL(c, hostId)+utils.GenerateEpsIDQuery(epsID), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r PremiumHostProtectStatus
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// deleteOpts whether to keep the premium domain policy when deleting the premium domain.
// Only used in the package.
type deleteOpts struct {
	KeepPolicy          bool   `q:"keepPolicy"`
	EnterpriseProjectID string `q:"enterprise_project_id"`
}

// Delete a premium domain by id.
func Delete(c *golangsdk.ServiceClient, hostId string, keepPolicy bool) (*SimplePremiumHost, error) {
	return DeleteWithEpsID(c, keepPolicy, hostId, "")
}

func DeleteWithEpsID(c *golangsdk.ServiceClient, keepPolicy bool, hostId, epsID string) (*SimplePremiumHost, error) {
	opts := deleteOpts{
		KeepPolicy:          keepPolicy,
		EnterpriseProjectID: epsID,
	}

	url := resourceURL(c, hostId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.DeleteWithResponse(url, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r SimplePremiumHost
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}
