package acls

import (
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure that used to create a new ACL policy.
type CreateOpts struct {
	// The ACL name.
	Name string `json:"acl_name" required:"true"`
	// The ACL type. The valid values are as follows:
	// + PERMIT
	// + DENY
	Type string `json:"acl_type" required:"true"`
	// The value of the ACL policy.
	// One or more values are supported, separated by commas (,).
	Value string `json:"acl_value" required:"true"`
	// The entity type. The valid values are as follows:
	// + IP
	// + DOMAIN
	// + DOMAIN_ID
	EntityType string `json:"entity_type" required:"true"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a private DNAT rule using given parameters.
func Create(c *golangsdk.ServiceClient, instanceId string, opts CreateOpts) (*Policy, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Policy
	_, err = c.Post(rootURL(c, instanceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method used to obtain the ACL policy detail by its ID.
func Get(c *golangsdk.ServiceClient, instanceId, policyId string) (*Policy, error) {
	var r Policy
	_, err := c.Get(resourceURL(c, instanceId, policyId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// UpdateOpts is the structure that used to modify an existing ACL policy configuration.
type UpdateOpts struct {
	// The ACL name.
	Name string `json:"acl_name" required:"true"`
	// The ACL type. The valid values are as follows:
	// + PERMIT
	// + DENY
	Type string `json:"acl_type" required:"true"`
	// The value of the ACL policy.
	// One or more values are supported, separated by commas (,).
	Value string `json:"acl_value" required:"true"`
	// The entity type. The valid values are as follows:
	// + IP
	// + DOMAIN
	// + DOMAIN_ID
	// The entity type does not support update.
	EntityType string `json:"entity_type" required:"true"`
}

// ListBindOpts is the structure used to querying published API list that ACL policy associated.
type ListOpts struct {
	// The instnace ID to which the API belongs.
	InstanceId string `json:"-" required:"true"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// The ACL policy ID.
	PolicyId string `q:"id"`
	// The ACL policy name.
	Name string `q:"name"`
	// The ACL type.
	Type string `q:"acl_type"`
	// The object type.
	EntityType string `q:"entity_type"`
	// Parameter name (name) for exact matching.
	PreciseSearch string `q:"precise_search"`
}

// List is a method to obtain all ACL policies under a specified instance.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Policy, error) {
	url := rootURL(c, opts.InstanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := PolicyPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractPolicies(pages)
}

// Update is a method used to modify the ACL policy configuration using given parameters.
func Update(c *golangsdk.ServiceClient, instanceId, policyId string, opts UpdateOpts) (*Policy, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Policy
	_, err = c.Put(resourceURL(c, instanceId, policyId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to remove the specified ACL policy using its ID and related dedicated instance ID.
func Delete(c *golangsdk.ServiceClient, instanceId, policyId string) error {
	_, err := c.Delete(resourceURL(c, instanceId, policyId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// BindOpts is the structure that used to bind a policy to the published APIs.
type BindOpts struct {
	// Instnace ID to which the API belongs.
	InstanceId string `json:"-" required:"true"`
	// ACL policy ID.
	PolicyId string `json:"acl_id,omitempty"`
	// The IDs of the API publish record.
	PublishIds []string `json:"publish_ids,omitempty"`
}

// Bind is a method to bind the policy to one or more APIs.
func Bind(c *golangsdk.ServiceClient, opts BindOpts) ([]BindResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r struct {
		BindList []BindResp `json:"acl_bindings"`
	}
	_, err = c.Post(bindURL(c, opts.InstanceId), b, &r, nil)
	return r.BindList, err
}

// ListBindOpts is the structure used to querying published API list that ACL policy associated.
type ListBindOpts struct {
	// The instnace ID to which the API belongs.
	InstanceId string `json:"-" required:"true"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// The ACL policy ID.
	PolicyId string `q:"acl_id"`
	// The API ID.
	ApiId string `q:"api_id"`
	// The API name.
	ApiName string `q:"api_name"`
	// The environment ID where the API is published.
	EnvId string `q:"env_id"`
	// The group ID where the API is located.
	GroupId string `q:"group_id"`
}

// ListBind is a method to obtain all API to which the ACL policy bound.
func ListBind(c *golangsdk.ServiceClient, opts ListBindOpts) ([]AclBindApiInfo, error) {
	url := listBindURL(c, opts.InstanceId)
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

// BatchUnbindOpts is the structure that used to unbinding policy from the published APIs.
type BatchUnbindOpts struct {
	// Instance ID.
	InstanceId string `json:"-" required:"true"`
	// List of API and ACL policy binding relationship IDs that need to be unbound.
	AclBindings []string `json:"acl_bindings,omitempty"`
}

// BatchUnbind is an API to unbind all ACL policies associated with the list.
func BatchUnbind(c *golangsdk.ServiceClient, unbindOpt BatchUnbindOpts, action string) (*BatchUnbindResp, error) {
	b, err := golangsdk.BuildRequestBody(unbindOpt, "")
	if err != nil {
		return nil, err
	}

	var (
		url = fmt.Sprintf("%v?action=%v", bindURL(c, unbindOpt.InstanceId), action)
		r   BatchUnbindResp
	)
	_, err = c.Put(url, b, &r, nil)
	return &r, err
}
