package components

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOpts is the structure required by the Create method to create a new component.
type CreateOpts struct {
	// Application component name.
	// The value can contain 2 to 64 characters, including letters, digits, hyphens (-), and underscores (_).
	// It must start with a letter and end with a letter or digit.
	Name         string `json:"name" required:"true"`
	WorkloadName string `json:"workload_name,omitempty"`
	// Description.
	// The value can contain up to 128 characters.
	Description string `json:"description,omitempty"`
	// Source of the code or software package.
	Source *Source `json:"source,omitempty" required:"true"`
	// Component builder.
	Build           *Build           `json:"build,omitempty"`
	Labels          []*KeyValue      `json:"labels,omitempty"`
	PodLabels       []*KeyValue      `json:"pod_labels,omitempty"`
	RuntimeStack    RuntimeStack     `json:"runtime_stack" required:"true"`
	LimitCpu        float64              `json:"limit_cpu,omitempty"`
	LimitMemory     float64              `json:"limit_memory,omitempty"`
	RequestCpu      float64              `json:"request_cpu,omitempty"`
	RequestMemory   float64              `json:"request_memory,omitempty"`
	Replica         int              `json:"replica"`
	Version         string           `json:"version" required:"true"`
	Envs            []*Env           `json:"envs,omitempty"`
	Storages        []*Storage       `json:"storage,omitempty"`
	DeployStrategy  *DeployStrategy  `json:"deploy_strategy,omitempty"`
	Command         *Command         `json:"command,omitempty"`
	PostStart       *K8sLifeCycle    `json:"post_start,omitempty"`
	PreStop         *K8sLifeCycle    `json:"pre_stop,omitempty"`
	Mesher          *Mesher          `json:"mesher,omitempty"`
	EnableSermantInjection bool      `json:"enable_sermant_injection,omitempty"`
	Timezone        string           `json:"timezone,omitempty"`
	JvmOpts         string           `json:"jvm_opts,omitempty"`
	TomcatOpts      *TomcatOpts      `json:"tomcat_opts,omitempty"`
	HostAliases     []*HostAlias     `json:"host_aliases,omitempty"`
	DnsPolicy       string           `json:"dns_policy,omitempty"`
	DnsConfig       *DnsConfig       `json:"dns_config,omitempty"`
	SecurityContext *SecurityContext `json:"security_context,omitempty"`
	WorkloadKind    string           `json:"workload_kind,omitempty"`
	Logs            []*Log            `json:"logs,omitempty"`
	CustomMetric    *CustomMetric    `json:"custom_metric,omitempty"`
	Affinity        *Affinity        `json:"affinity,omitempty"`
	AntiAffinity    *Affinity        `json:"anti_affinity,omitempty"`
	LivenessProbe   *K8sProbe        `json:"liveness_probe,omitempty"`
	ReadinessProbe  *K8sProbe        `json:"readiness_probe,omitempty"`
	ReferResources  []*Resource      `json:"refer_resources,omitempty"`
	// Environment info
	EnvironmentID string `json:"environment_id" required:"true"`
	// Application info
	ApplicationID string `json:"application_id" required:"true"`
	// The enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id"`
}

// Create is a method to create a component in the specified application using given parameters.
func Create(c *golangsdk.ServiceClient, appId string, opts CreateOpts) (*JobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r JobResp
	_, err = c.Post(rootURL(c, appId), b, &r, nil)
	return &r, err
}

// Get is a method to retrieves a particular configuration based on its unique ID.
func Get(c *golangsdk.ServiceClient, appId, componentId string) (*Component, error) {
	var r Component
	_, err := c.Get(resourceURL(c, appId, componentId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// Value range: 0–100, or 1000.
	// Default value: 1000, indicating that a maximum of 1000 records can be queried and all records are displayed on
	// the same page.
	Limit string `q:"limit"`
	// The offset number.
	Offset int `q:"offset"`
	// Sorting field. By default, query results are sorted by creation time.
	// The following enumerated values are supported: create_time, name, and update_time.
	OrderBy string `q:"order_by"`
	// Descending or ascending order. Default value: desc.
	Order string `q:"order"`
}

// List is a method to query the list of the components using given opts.
func List(c *golangsdk.ServiceClient, appId string, opts ListOpts) ([]Component, error) {
	url := rootURL(c, appId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := ComponentPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractComponents(pages)
}

// UpdateOpts is the structure required by the Update method to update the component configuration.
type UpdateOpts struct {
	// Application component name.
	// The value can contain 2 to 64 characters, including letters, digits, hyphens (-), and underscores (_).
	// It must start with a letter and end with a letter or digit.
	Name string `json:"name,omitempty"`
	// Description.
	// The value can contain up to 128 characters.
	Description *string `json:"description,omitempty"`
	// Source of the code or software package.
	Source *Source `json:"source,omitempty"`
	// Component build.
	Builder *Build `json:"build,omitempty"`
}

// Update is a method to update the component configuration, such as name, description, builder and code source.
func Update(c *golangsdk.ServiceClient, appId, componentId string, opts UpdateOpts) (*Component, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Component
	_, err = c.Put(resourceURL(c, appId, componentId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to delete an existing component from a specified application.
func Delete(c *golangsdk.ServiceClient, appId, componentId string) error {
	_, err := c.Delete(resourceURL(c, appId, componentId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
