package plugins

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure used to create a new plugin.
type CreateOpts struct {
	// The ID of the dedicated instance to which the plugin belongs.
	InstanceId string `json:"-" required:"true"`
	// Plugin name.
	// The valid length is limited from `3` to `255`, only Chinese characters, English letters, digits, hyphens (-) and
	// underscores (_) are allowed. The name must start with an English letter or Chinese character.
	Name string `json:"plugin_name" required:"true"`
	// Plugin type.
	// The valid values are as follows:
	// + 'cors': CORS, specify preflight request headers and response headers and automatically create preflight
	//   request APIs for cross-origin API access.
	// + 'set_resp_headers': HTTP Response Header Management, customize HTTP headers that will be contained in an API
	//   response.
	// + 'rate_limit': Request Throttling 2.0, limits the number of times an API can be called within a specific time
	//   period. It supports parameter-based, basic, and excluded throttling.
	// + 'kafka_log': Kafka Log Push, Push detailed API calling logs to kafka for you to easily obtain logs.
	// + 'breaker': Circuit Breaker, circuit breaker protect the system when performance issues occur on backend
	//   service.
	Type string `json:"plugin_type" required:"true"`
	// The available scope for plugin, the valid value is 'global'.
	Scope string `json:"plugin_scope" required:"true"`
	// The configuration details for plugin.
	Content string `json:"plugin_content" required:"true"`
	// The plugin description.
	// The valid length is limited from `3` to `255` characters.
	Description string `json:"remark,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new plugin using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Plugin, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Plugin
	_, err = c.Post(rootURL(c, opts.InstanceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to obtain an existing plugin detail by its ID and related instance ID.
func Get(client *golangsdk.ServiceClient, instanceId, pluginId string) (*Plugin, error) {
	var r Plugin
	_, err := client.Get(resourceURL(client, instanceId, pluginId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// UpdateOpts is the structure used to update the plugin configuration.
type UpdateOpts struct {
	// The instnace ID to which the plugin belongs.
	InstanceId string `json:"-" required:"true"`
	// The plugin ID.
	ID string `json:"-" required:"true"`
	// Plugin name.
	// The valid length is limited from `3` to `255`, only Chinese characters, English letters, digits, hyphens (-) and
	// underscores (_) are allowed. The name must start with an English letter or Chinese character.
	Name string `json:"plugin_name" required:"true"`
	// Plugin type.
	// The valid values are as follows:
	// + 'cors': CORS, specify preflight request headers and response headers and automatically create preflight
	//   request APIs for cross-origin API access.
	// + 'set_resp_headers': HTTP Response Header Management, customize HTTP headers that will be contained in an API
	//   response.
	// + 'rate_limit': Request Throttling 2.0, limits the number of times an API can be called within a specific time
	//   period. It supports parameter-based, basic, and excluded throttling.
	// + 'kafka_log': Kafka Log Push, Push detailed API calling logs to kafka for you to easily obtain logs.
	// + 'breaker': Circuit Breaker, circuit breaker protect the system when performance issues occur on backend
	//   service.
	Type string `json:"plugin_type" required:"true"`
	// The available scope for plugin, the valid value is 'global'.
	Scope string `json:"plugin_scope" required:"true"`
	// The configuration details for plugin.
	Content string `json:"plugin_content" required:"true"`
	// The plugin description.
	// The valid length is limited from `3` to `255` characters.
	Description *string `json:"remark,omitempty"`
}

// Update is a method used to update a plugin using given parameters.
func Update(c *golangsdk.ServiceClient, opts UpdateOpts) (*Plugin, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Plugin
	_, err = c.Put(resourceURL(c, opts.InstanceId, opts.ID), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to remove the specified plugin using its ID and related dedicated instance ID.
func Delete(c *golangsdk.ServiceClient, instanceId, pluginId string) error {
	_, err := c.Delete(resourceURL(c, instanceId, pluginId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// BindOpts is the structure that used to bind a plugin to the published APIs.
type BindOpts struct {
	// The instnace ID to which the plugin belongs.
	InstanceId string `json:"-" required:"true"`
	// The plugin ID.
	PluginId string `json:"-" required:"true"`
	// The environment ID where the API is published.
	EnvId string `json:"env_id" required:"true"`
	// The IDs of the API publish record.
	ApiIds []string `json:"api_ids" required:"true"`
}

// Bind is a method to bind a plugin to one or more APIs.
func Bind(c *golangsdk.ServiceClient, opts BindOpts) ([]PluginBindDetail, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r BindResp
	_, err = c.Post(bindURL(c, opts.InstanceId, opts.PluginId), b, &r, nil)
	return r.Bindings, err
}

// ListBindOpts is the structure used to querying published API list that plugin associated.
type ListBindOpts struct {
	// The instnace ID to which the plugin belongs.
	InstanceId string `json:"-" required:"true"`
	// The plugin ID.
	PluginId string `json:"-" required:"true"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// The environment ID where the API is published.
	EnvId string `q:"env_id"`
	// The API name.
	ApiName string `q:"api_name"`
	// The API ID.
	ApiId string `q:"api_id"`
	// The group ID where the API is located.
	GroupId string `q:"group_id"`
	// The request method.
	RequestMethod string `q:"req_method"`
	// The request URI.
	RequestUri string `q:"req_uri"`
}

// ListBind is a method to obtain all API to which the plugin bound.
func ListBind(c *golangsdk.ServiceClient, opts ListBindOpts) ([]BindApiInfo, error) {
	url := listBindURL(c, opts.InstanceId, opts.PluginId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := BindPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractBindInfos(pages)
}

// UnbindOpts is the structure that used to unbind the published APIs from the plugin.
type UnbindOpts struct {
	// The instnace ID to which the plugin belongs.
	InstanceId string `json:"-" required:"true"`
	// The plugin ID.
	PluginId string `json:"-" required:"true"`
	// The environment ID where the API is published.
	EnvId string `json:"env_id" required:"true"`
	// The IDs of the API publish record.
	ApiIds []string `json:"api_ids" required:"true"`
}

// Unbind is an method used to unbind one or more APIs from the plugin.
func Unbind(c *golangsdk.ServiceClient, opts UnbindOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(unbindURL(c, opts.InstanceId, opts.PluginId), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
		OkCodes:     []int{204},
	})
	return err
}
