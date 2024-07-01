package plugins

import "github.com/chnsz/golangsdk/pagination"

// Plugin is the structure that represents the plugin detail.
type Plugin struct {
	// The plugin ID.
	ID string `json:"plugin_id"`
	// Plugin name.
	// The valid length is limited from `3` to `255`, only Chinese characters, English letters, digits, hyphens (-) and
	// underscores (_) are allowed. The name must start with an English letter or Chinese character.
	Name string `json:"plugin_name"`
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
	Type string `json:"plugin_type"`
	// The available scope for plugin, the valid value is 'global'.
	Scope string `json:"plugin_scope"`
	// The configuration details for plugin.
	Content string `json:"plugin_content"`
	// The plugin description.
	// The valid length is limited from `3` to `255` characters.
	Description string `json:"remark"`
	// The creation time.
	CreatedAt string `json:"create_time"`
	// The latest update time.
	UpdatedAt string `json:"update_time"`
}

// BindResp is the structure that represents the API response of the plugin binding.
type BindResp struct {
	// The published APIs of the binding relationship.
	Bindings []PluginBindDetail `json:"bindings"`
}

// PluginBindDetail is the structure that represents the binding details.
type PluginBindDetail struct {
	// The binding ID.
	BindId string `json:"plugin_attach_id"`
	// The plugin ID.
	PluginId string `json:"plugin_id"`
	// The plugin name.
	PluginName string `json:"plugin_name"`
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
	PluginType string `json:"plugin_type"`
	// The available scope for plugin, the valid value is 'global'.
	PluginScope string `json:"plugin_scope"`
	// The environment ID where the API is published.
	EnvId string `json:"env_id"`
	// The name of the environment published by the API.
	EnvName string `json:"env_name"`
	// The API ID.
	ApiId string `json:"api_id"`
	// The API name.
	ApiName string `json:"apig_name"`
	// The time when the API and the plugin were bound.
	BoundAt string `json:"attached_time"`
}

// BindPage is a single page maximum result representing a query by offset page.
type BindPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a BindPage struct is empty.
func (b BindPage) IsEmpty() (bool, error) {
	arr, err := ExtractBindInfos(b)
	return len(arr) == 0, err
}

// ExtractBindInfos is a method to extract the list of binding details for plugin.
func ExtractBindInfos(r pagination.Page) ([]BindApiInfo, error) {
	var s []BindApiInfo
	err := r.(BindPage).Result.ExtractIntoSlicePtr(&s, "apis")
	return s, err
}

// BindApiInfo is an object that represents the bind detail.
type BindApiInfo struct {
	// API ID.
	ApiId string `json:"api_id"`
	// API name.
	ApiName string `json:"api_name"`
	// API type.
	Type int `json:"type"`
	// The request protocol of the API.
	RequestProtocol string `json:"req_protocol"`
	// The request method of the API.
	RequestMethod string `json:"req_method"`
	// The request URI of the API.
	RequestUri string `json:"request_uri"`
	// The authorization type of the API.
	AuthType string `json:"auth_type"`
	// The match mode of the API.
	MatchMode string `json:"match_mode"`
	// The description of the API.
	Remark string `json:"remark"`
	// The ID of API group to which the API belongs.
	GroupId string `json:"group_id"`
	// The name of API group to which the API belongs.
	GroupName string `json:"group_name"`
	// Home integration application code, which is compatible with the field of the Roma instance.
	// Generally, the value is null.
	RomaAppId string `json:"roma_app_id"`
	// The ID of environment where the API was published.
	EnvId string `json:"env_id"`
	// The name of environment where the API was published.
	EnvName string `json:"env_name"`
	// The API publish ID.
	PublishId string `json:"publish_id"`
	// The ID of plugin policy binding.
	BindId string `json:"plugin_attach_id"`
	// The bound time.
	BoundAt string `json:"attached_time"`
}
