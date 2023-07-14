package clouds

import (
	"fmt"
	"github.com/chnsz/golangsdk"
)

// CreateOpts is the structure required by the 'Create' method to create a cloud WAF.
type CreateOpts struct {
	// The ID of the project to which the cloud WAF belongs.
	ProjectId string `json:"project_id" required:"true"`
	// Whether the order is auto pay.
	IsAutoPay *bool `json:"is_auto_pay" required:"true"`
	// Whether auto renew is enabled for resource payment.
	IsAutoRenew *bool `json:"is_auto_renew" required:"true"`
	// The region where the cloud WAF is located.
	RegionId string `json:"region_id" required:"true"`
	// The configuration of the cloud WAF, such as specification code.
	ProductInfo *ProductInfo `json:"waf_product_info,omitempty"`
	// The configuration of the bandwidth extended packages.
	BandwidthExpackProductInfo *ExpackProductInfo `json:"bandwidth_expack_product_info,omitempty"`
	// The configuration of the domain extended packages.
	DomainExpackProductInfo *ExpackProductInfo `json:"domain_expack_product_info,omitempty"`
	// The configuration of the rule extended packages.
	RuleExpackProductInfo *ExpackProductInfo `json:"rule_expack_product_info,omitempty"`
	// The ID of the enterprise project to which the cloud WAF belongs.
	EnterpriseProjectId string `q:"enterprise_project_id" json:"-"`
}

// ProductInfo is an object that represents the configuration of the cloud WAF.
type ProductInfo struct {
	// The specification of the cloud WAF.
	ResourceSpecCode string `json:"resource_spec_code,omitempty"`
	// The charging period unit of the cloud WAF.
	PeriodType string `json:"period_type,omitempty"`
	// The charging period of the cloud WAF.
	PeriodNum int `json:"period_num,omitempty"`
}

// ExpackProductInfo is an object that represents the configuration of the extended packages.
type ExpackProductInfo struct {
	// The number of the extended packages
	ResourceSize int `json:"resource_size,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new cloud WAF using given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (*string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	url := createURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r createResp
	_, err = client.Post(url, b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.OrderId, err
}

// Get is a method used to obtain the cloud WAF details.
func Get(client *golangsdk.ServiceClient) (*Instance, error) {
	return GetWithEpsID(client, "")
}

// GetWithEpsID is a method used to obtain the cloud WAF details with eps ID.
func GetWithEpsID(client *golangsdk.ServiceClient, epsID string) (*Instance, error) {
	var r Instance
	_, err := client.Get(getURL(client)+generateEpsIdQuery(epsID), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

func generateEpsIdQuery(epsID string) string {
	if len(epsID) == 0 {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsID)
}

// UpdateOpts is the structure required by the 'Update' method to update the cloud WAF configuration.
type UpdateOpts struct {
	// The ID of the project to which the cloud WAF belongs.
	ProjectId string `json:"project_id" required:"true"`
	// Whether the order is auto pay.
	IsAutoPay *bool `json:"is_auto_pay" required:"true"`
	// The configuration of the cloud WAF, such as specification code.
	ProductInfo *UpdateProductInfo `json:"waf_product_info,omitempty"`
	// The configuration of the bandwidth extended packages.
	BandwidthExpackProductInfo *ExpackProductInfo `json:"bandwidth_expack_product_info,omitempty"`
	// The configuration of the domain extended packages.
	DomainExpackProductInfo *ExpackProductInfo `json:"domain_expack_product_info,omitempty"`
	// The configuration of the rule extended packages.
	RuleExpackProductInfo *ExpackProductInfo `json:"rule_expack_product_info,omitempty"`
	// The ID of the enterprise project to which the cloud WAF belongs.
	EnterpriseProjectId string `q:"enterprise_project_id" json:"-"`
}

// UpdateProductInfo is an object that represents the update configuration of the cloud WAF.
type UpdateProductInfo struct {
	// Whether the AS path attributes of the routes are not compared during load balancing.
	ResourceSpecCode string `json:"resource_spec_code,omitempty"`
}

// Update is a method used to update the cloud WAF using given parameters.
func Update(client *golangsdk.ServiceClient, opts UpdateOpts) (*string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	url := updateURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r updateResp
	_, err = client.Post(url, b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.OrderId, err
}

// CreatePostPaidOpts is the structure required by the 'Create' method to create a post paid cloud WAF.
type CreatePostPaidOpts struct {
	// The region where the cloud WAF is located. This field will be set to header
	Region string `json:"-" required:"true"`
	// The website to which the account belongs.
	ConsoleArea string `json:"console_area" required:"true"`
	// The ID of the enterprise project to which the cloud WAF belongs.
	EnterpriseProjectId string `q:"enterprise_project_id" json:"-"`
}

// CreatePostPaid is a method used to create a new post paid cloud WAF using given parameters.
func CreatePostPaid(client *golangsdk.ServiceClient, opts CreatePostPaidOpts) (*Instance, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	url := createOrDeletePostPaidURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()
	var r Instance

	moreHeaders := requestOpts.MoreHeaders
	moreHeaders["region"] = opts.Region
	_, err = client.Post(url, b, &r, &golangsdk.RequestOpts{
		MoreHeaders: moreHeaders,
	})
	return &r, err
}

// DeletePostPaidOpts is the structure required by the 'Delete' method to delete a post paid cloud WAF.
type DeletePostPaidOpts struct {
	// The region where the cloud WAF is located. This field will be set to header
	Region string `json:"-" required:"true"`
	// The ID of the enterprise project to which the cloud WAF belongs.
	EnterpriseProjectId string `q:"enterprise_project_id" json:"-"`
}

// DeletePostPaid is a method used to delete a post paid cloud WAF using given parameters.
func DeletePostPaid(client *golangsdk.ServiceClient, opts DeletePostPaidOpts) error {
	url := createOrDeletePostPaidURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return err
	}
	url += query.String()
	moreHeaders := requestOpts.MoreHeaders
	moreHeaders["region"] = opts.Region
	_, err = client.Delete(url, &golangsdk.RequestOpts{
		MoreHeaders: moreHeaders,
	})
	return err
}
