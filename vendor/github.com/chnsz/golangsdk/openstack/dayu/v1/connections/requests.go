package connections

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure that used by 'Create' method to create a new data connection.
type CreateOpts struct {
	// The ID of the workspace where the data connection is located.
	WorkspaceId string `json:"-" required:"true"`
	// The list of structure for data source configuration.
	DataSourceVos []DataSourceVo `json:"data_source_vos" required:"true"`
}

// DataSourceVo is the structure that represents the configuration of the data source.
type DataSourceVo struct {
	// The data connection name.
	DwName string `json:"dw_name" required:"true"`
	// The data connection type.
	DwType string `json:"dw_type" required:"true"`
	// The dynamic configuration for the specified type of data source.
	DwConfig interface{} `json:"dw_config" required:"true"`
	// The agent ID.
	AgentId string `json:"agent_id,omitempty"`
	// The agent name.
	AgentName string `json:"agent_name,omitempty"`
	// The data connection mode.
	EnvType int `json:"env_type,omitempty"`
}

func buildRequestOpts(workspaceId string) *golangsdk.RequestOpts {
	return &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
			"workspace":    workspaceId, // workspace is required.
		},
	}
}

// Create is a method that used to create a new connection using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return "", err
	}

	var r createResp
	_, err = c.Post(rootURL(c), b, &r, buildRequestOpts(opts.WorkspaceId))
	return r.DataConnectionId, err
}

// ValidateOpts is the structure that represents the pre-check configuration.
type ValidateOpts struct {
	// The ID of the workspace where the data connection is located.
	WorkspaceId string `json:"-" required:"true"`
	// The data connection name.
	DwName string `json:"dw_name" required:"true"`
	// The data connection type.
	DwType string `json:"dw_type" required:"true"`
	// The dynamic configuration for the specified type of data source.
	DwConfig interface{} `json:"dw_config" required:"true"`
	// The agent ID.
	AgentId string `json:"agent_id,omitempty"`
	// The agent name.
	AgentName string `json:"agent_name,omitempty"`
	// The data connection mode.
	EnvType int `json:"env_type,omitempty"`
}

// Validate is a method that used to pre-check the configuration of the data connection using given parameters.
func Validate(c *golangsdk.ServiceClient, opts ValidateOpts) (*ValidateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r ValidateResp
	_, err = c.Post(validateURL(c), b, &r, buildRequestOpts(opts.WorkspaceId))
	return &r, err
}

// Get is a method to obtain the data connection detail by its ID and related workspace ID.
func Get(c *golangsdk.ServiceClient, workspaceId, connectionId string) (*Connection, error) {
	var r Connection
	_, err := c.Get(resourceURL(c, connectionId), &r, buildRequestOpts(workspaceId))
	return &r, err
}

// ListOpts is the structure used by 'List' method to query data connections.
type ListOpts struct {
	// The ID of the workspace where the data connection is located.
	WorkspaceId string `json:"-" required:"true"`
	// The data connection name.
	Name string `q:"name"`
	// The data connection type.
	Type string `q:"type"`
	// Limit is the records count to be queried.
	Limit string `q:"limit"`
	// Offset number.
	Offset string `q:"offset"`
}

// List is a method to query the list of the data connections using given parameters.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Connection, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pager := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ConnectionPage{pagination.OffsetPageBase{PageResult: r}}
	})
	queryOpts := buildRequestOpts(opts.WorkspaceId)
	pager.Headers = queryOpts.MoreHeaders
	pages, err := pager.AllPages()

	if err != nil {
		return nil, err
	}
	return extractConnections(pages)
}

// UpdateOpts is the structure that used by 'Update' method to modify the configuration of the data connection.
type UpdateOpts struct {
	// The connection ID.
	ConnectionId string `json:"-" required:"true"`
	// The ID of the workspace where the data connection is located.
	WorkspaceId string `json:"-" required:"true"`
	// The list of structure for data source configuration.
	DataSourceVos []DataSourceVo `json:"data_source_vos" required:"true"`
}

// Update is a method to modify the specified connection using given parameters.
func Update(c *golangsdk.ServiceClient, opts UpdateOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(resourceURL(c, opts.ConnectionId), b, nil, buildRequestOpts(opts.WorkspaceId))
	return err
}

// Delete is a method to delete a specified connection by its ID and the related workspace ID.
func Delete(client *golangsdk.ServiceClient, workspaceId, connnectionId string) error {
	_, err := client.Delete(resourceURL(client, connnectionId), buildRequestOpts(workspaceId))
	return err
}
