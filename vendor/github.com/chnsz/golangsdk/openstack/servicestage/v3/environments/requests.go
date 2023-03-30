package environments

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the Create method to create a new environment.
type CreateOpts struct {
	// Specified the environment name with 2 to 64 characters long.
	// It consists of English letters, numbers, underscores (-), and underscores (_).
	// It must start with an English letter and end with an English letter or number.
	Name string `json:"name" required:"true"`
	// Environment deploy mode. Value: container, virtualmachine or mixed.
	DeployMode string `json:"deploy_mode"`
	// Specified the environment description.
	// The description can contain a maximum of 96 characters.
	Description string `json:"description,omitempty"`
	// Specified the enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// Specifies the VPC ID.
	VpcId string `json:"vpc_id" required:"true"`
	// Environment's Labels
	Labels []Label `json:"labels"`
	// It only takes effect when the environment's DeployMode is "virtualmachine"
	// Value: 50 or 500
	VmClusterSize int `json:"vm_cluster_size"`
	// Specified the environment alias.
	// The alias can contain a maximum of 96 characters.
	Alias string `json:"alias,omitempty"`
	// Specified the billing mode. The valid values are:
	//   provided: provided resources are used and no fees are charged.
	//   on_demanded: on-demand charging.
	//   monthly: monthly subscription.
	ChargeMode string `json:"charge_mode,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a new ServiceStage environment using create option.
// Environment is a collection of infrastructures, covering computing, storage and networks, used for application
// deployment and running.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Environment, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst Environment
	_, err = c.Post(rootURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &rst, err
}

// Get is a method to obtain the details of a specified ServiceStage environment using its ID.
func Get(c *golangsdk.ServiceClient, envId string) (*Environment, error) {
	var r Environment
	_, err := c.Get(detailURL(c, envId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// Value range: 0–100.
	// Default value: 1000, indicating that a maximum of 1000 records can be queried and all records are displayed on
	// the same page.
	Limit int `q:"limit"`
	// The offset number.
	Offset int `q:"offset"`
	// Sorting field. By default, query results are sorted by creation time.
	// The following enumerated values are supported: create_time, name, and update_time.
	OrderBy string `q:"order_by"`
	// Descending or ascending order. Default value: desc.
	Order string `q:"order"`
}

// List is a method to query the list of the environments using given opts.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Environment, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := EnvironmentPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractEnvironments(pages)
}

// UpdateOpts is the structure required by the Update method to update the configuration of the specified environment.
type UpdateOpts struct {
	// Specified the environment description.
	// The description can contain a maximum of 128 characters.
	Description string `json:"description,omitempty"`
	// Specified the environment name with 2 to 64 characters long.
	// It consists of English letters, numbers, underscores (-), and underscores (_).
	// It must start with an English letter and end with an English letter or number.
	Name string `json:"name,omitempty"`
	// Specified the environment alias.
	// The alias can contain a maximum of 64 characters.
	Alias string `json:"alias,omitempty"`
	// The enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// Environment's Labels
	Labels []Label `json:"labels"`
}

// Update is a method to update the current dependency configuration.
func Update(c *golangsdk.ServiceClient, envId string, opts UpdateOpts) (*Environment, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Environment
	_, err = c.Put(detailURL(c, envId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to remove an existing environment.
func Delete(c *golangsdk.ServiceClient, envId string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(detailURL(c, envId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r
}

// ListResources is a method to query the list of resources from the specified environment.
func ListResources(c *golangsdk.ServiceClient, envId string) ([]Resource, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, envId), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	var r []Resource
	err = rst.ExtractIntoSlicePtr(&r, "resources")
	return r, err
}

// ResourceOpts is a structure required but the UpdateResources method to bind or unbind the resources.
type ResourceOpts struct {
	// Basic resources to be added.
	// VPC, EIP, ELB, CCE, CCI, ECS, AS, CSE, DCS, RDS, PVC
	Resources []Resource `json:"resources,omitempty"`
}

// UpdateResources is a method to add or remove the basic resources and the optional resources.
func UpdateResources(c *golangsdk.ServiceClient, envId string, opts ResourceOpts) ([]Resource, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, envId), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	var r []Resource
	err = rst.ExtractIntoSlicePtr(&r, "resources")
	return r, err
}
