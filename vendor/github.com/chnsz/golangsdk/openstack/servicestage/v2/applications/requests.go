package applications

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the Create method to create a new application.
type CreateOpts struct {
	// Specified the application name with 2 to 64 characters long.
	// It consists of English letters, numbers, underscores (-), and underscores (_).
	// It must start with an English letter and end with an English letter or number.
	Name string `json:"name" required:"true"`
	// Specified the application description.
	// The description can contain a maximum of 96 characters.
	Description *string `json:"description,omitempty"`
	// Specified the enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a new ServiceStage application using create option.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Application, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst Application
	_, err = c.Post(rootURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &rst, err
}

// Get is a method to obtain the details of a specified ServiceStage application using its ID.
func Get(c *golangsdk.ServiceClient, appId string) (*Application, error) {
	var r Application
	_, err := c.Get(resourceURL(c, appId), &r, &golangsdk.RequestOpts{
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
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Application, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := ApplicationPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractApplications(pages)
}

// UpdateOpts is the structure required by the Update method to update the configuration of the specified environment.
type UpdateOpts struct {
	// Specified the application description.
	// The description can contain a maximum of 96 characters.
	Description *string `json:"description,omitempty"`
	// Specified the environment name with 2 to 64 characters long.
	// It consists of English letters, numbers, underscores (-), and underscores (_).
	// It must start with an English letter and end with an English letter or number.
	Name string `json:"name,omitempty"`
}

// Update is a method to update the current dependency configuration.
func Update(c *golangsdk.ServiceClient, appId string, opts UpdateOpts) (*Application, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Application
	_, err = c.Put(resourceURL(c, appId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to remove an existing application.
func Delete(c *golangsdk.ServiceClient, appId string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(resourceURL(c, appId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r
}

// ConfigOpts is a structure required bu the UpdateConfig method to bind or unbind the resources.
type ConfigOpts struct {
	// The environment ID.
	EnvId string `json:"environment_id" required:"true"`
	// The application configurations, such as public environment variables.
	Configuration Configuration `json:"configuration" required:"true"`
}

// Configuration is an object specifying the configuration of the application for a specified environment.
type Configuration struct {
	// The environment variables of the application.
	// If the names of multiple environment variables are the same, only the last environment variable takes effact.
	EnvVariables []Variable `json:"env,omitempty"`
}

// Variable is an object specifying the key/value pair of the environment variable.
type Variable struct {
	// The name of the environment variable.
	// The name contains 1 to 64 characters, including letters, digits, underscores (_), hyphens (-) and dots (.), and
	// cannot start with a digit or dot.
	Name string `json:"name" required:"true"`
	// The value of the environment variables.
	Value string `json:"value" required:"true"`
}

// UpdateConfig is a method to add or remove the basic resources and the optional resources.
func UpdateConfig(c *golangsdk.ServiceClient, appId, envId string, config Configuration) (*Application, error) {
	opts := ConfigOpts{
		EnvId:         envId,
		Configuration: config,
	}
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Application
	_, err = c.Put(configURL(c, appId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListConfigOpts allows to filter list configurations using given parameters.
type ListConfigOpts struct {
	// The environment ID. If this parameter is not specified, all environment are queried.
	EnvId string `q:"environment_id"`
}

// ListConfig is a method to query the list of configurations from the specified application.
func ListConfig(c *golangsdk.ServiceClient, appId string, opts ListConfigOpts) ([]ConfigResp, error) {
	url := configURL(c, appId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.Get(url, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err != nil {
		return nil, err
	}
	var r []ConfigResp
	rst.ExtractIntoSlicePtr(&r, "configuration")
	return r, nil
}

// DeleteConfig is a method to remove an existing configuration or some variables from a specified environment.
func DeleteConfig(c *golangsdk.ServiceClient, appId, envId string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult

	url := configURL(c, appId)
	opts := ListConfigOpts{
		EnvId: envId,
	}
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		r.Err = err
	}
	url += query.String()

	_, r.Err = c.Delete(url, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r
}
