package instances

import (
	"encoding/json"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new instance.
type CreateOpts struct {
	Region              string `json:"region_id" required:"true"`
	AvailabilityZone    string `json:"availability_zone" required:"true"`
	Name                string `json:"instance_name" required:"true"`
	SpecCode            string `json:"resource_spec_code" required:"true"`
	VpcID               string `json:"vpc_id" required:"true"`
	SubnetID            string `json:"net_id" required:"true"`
	SecurityGroupID     string `json:"security_group_id" required:"true"`
	EnterpriseProjectID string `json:"eps_id,omitempty"`

	// 2 - monthly; 3 - yearly
	PeriodType int `json:"period_type" required:"true"`
	// monthly: 1-9; yearly: 1-3
	PeriodNum int `json:"period_num" required:"true"`
	// 0 - false; 1 - true
	IsAutoRenew *int `json:"is_auto_renew" required:"true"`

	Tags []tags.ResourceTag `json:"tags,omitempty"`

	// the following parameters are reserved
	ProductID         string `json:"product_id,omitempty"`
	CommodityID       string `json:"commodity_id,omitempty"`
	PromotionInfo     string `json:"promotion_info,omitempty"`
	ExtPackageType    string `json:"extesion_package_type,omitempty"`
	BindingInstanceID string `json:"binding_instance_id,omitempty"`
	CdmVersion        string `json:"cdm_version,omitempty"`
	CloudServiceType  string `json:"cloud_service_type,omitempty"`
	ResourceType      string `json:"resource_type,omitempty"`
}

// ToInstanceCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToInstanceCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new instance.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(createURL(c), b, &r.Body, reqOpt)
	return
}

// ListOpts allows the filtering collections through the API.
type ListOpts struct {
	// Limit is the records count to be queried, the value ranges 1-100, the default value is 20
	Limit int `q:"limit"`

	// Offset number, the default value is 0
	Offset int `q:"offset"`
}

// ToInstanceListQuery formats a ListOpts into a query string
func (opts ListOpts) ToInstanceListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns collection of instances.
func List(c *golangsdk.ServiceClient, opts *ListOpts) ([]Instance, error) {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToInstanceListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}

	var s struct {
		Instances []Instance `json:"commodity_orders"`
	}
	_, err := c.Get(url, &s, nil)
	if err != nil {
		return nil, err
	}

	return s.Instances, nil
}

// ListAll returns all of instances start with initURL.
func ListAll(c *golangsdk.ServiceClient, opts *ListOpts) ([]Instance, error) {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToInstanceListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}

	resp, err := pagination.ListAllItems(c, pagination.Offset, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var s struct {
		Instances []Instance `json:"commodity_orders"`
	}
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, err
	}

	return s.Instances, nil
}
