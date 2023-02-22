package policyassignments

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the 'Create' method to assign a specified policy.
type CreateOpts struct {
	// The name of the policy assignment.
	Name string `json:"name" required:"true"`
	// The description of the policy assignment.
	Description string `json:"description,omitempty"`
	// The type of the policy assignment.
	// The valid values are as follow:
	// + builtin
	// + custom
	Type string `json:"policy_assignment_type,omitempty"`
	// The period of the policy rule check.
	Period string `json:"period,omitempty"`
	// The configuration used to filter resources.
	PolicyFilter PolicyFilter `json:"policy_filter,omitempty"`
	// The ID of the policy definition.
	PolicyDefinitionId string `json:"policy_definition_id,omitempty"`
	// The configuration of the custom policy.
	CustomPolicy *CustomPolicy `json:"custom_policy,omitempty"`
	// The rule definition of the policy assignment.
	Parameters map[string]PolicyParameterValue `json:"parameters,omitempty"`
}

// PolicyFilter is an object specifying the filter parameters.
type PolicyFilter struct {
	// The name of the region to which the filtered resources belong.
	RegionId string `json:"region_id,omitempty"`
	// The service name to which the filtered resources belong.
	ResourceProvider string `json:"resource_provider,omitempty"`
	// The resource type of the filtered resources.
	ResourceType string `json:"resource_type,omitempty"`
	// The ID used to filter resources.
	ResourceId string `json:"resource_id,omitempty"`
	// The tag name used to filter resources.
	TagKey string `json:"tag_key,omitempty"`
	// The tag value used to filter resources.
	TagValue string `json:"tag_value,omitempty"`
}

// CustomPolicy is an object specifying the configuration of the custom policy.
type CustomPolicy struct {
	// The function URN used to create the custom policy.
	FunctionUrn string `json:"function_urn" required:"true"`
	// The authorization type of the custom policy.
	AuthType string `json:"auth_type" required:"true"`
	// The authorization value of the custom policy.
	AuthValue map[string]interface{} `json:"auth_value,omitempty"`
}

// PolicyParameterValue is an object specifying the definition of the policy parameter value.
type PolicyParameterValue struct {
	// The value of the rule definition.
	Value interface{} `json:"value,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to assign a policy using given parameters.
func Create(c *golangsdk.ServiceClient, domainId string, opts CreateOpts) (*Assignment, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Assignment
	_, err = c.Put(rootURL(c, domainId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to obtain the assignment detail by its ID.
func Get(c *golangsdk.ServiceClient, domainId, assignmentId string) (*Assignment, error) {
	var r Assignment
	_, err := c.Get(resourceURL(c, domainId, assignmentId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// UpdateOpts is the structure required by the Update method to update the assignment configuration.
type UpdateOpts struct {
	// The name of the policy assignment.
	Name string `json:"name" required:"true"`
	// The description of the policy assignment.
	Description *string `json:"description,omitempty"`
	// The type of the policy assignment.
	// The valid values are as follow:
	// + builtin
	// + custom
	Type string `json:"policy_assignment_type,omitempty"`
	// The period of the policy rule check.
	Period string `json:"period,omitempty"`
	// The configuration used to filter resources.
	PolicyFilter PolicyFilter `json:"policy_filter,omitempty"`
	// The ID of the policy definition.
	PolicyDefinitionId string `json:"policy_definition_id,omitempty"`
	// The configuration of the custom policy.
	CustomPolicy *CustomPolicy `json:"custom_policy,omitempty"`
	// The rule definition of the policy assignment.
	Parameters map[string]PolicyParameterValue `json:"parameters,omitempty"`
}

// Update is a method to update the assignment configuration.
func Update(c *golangsdk.ServiceClient, domainId, assignmentId string, opts UpdateOpts) (*Assignment, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Assignment
	_, err = c.Put(resourceURL(c, domainId, assignmentId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Enable is a method to enable the resource function of the policy assignment.
func Enable(c *golangsdk.ServiceClient, domainId, assignmentId string) error {
	_, err := c.Post(enableURL(c, domainId, assignmentId), nil, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// Disable is a method to disable the resource function of the policy assignment.
func Disable(c *golangsdk.ServiceClient, domainId, assignmentId string) error {
	_, err := c.Post(disableURL(c, domainId, assignmentId), nil, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// Delete is a method to remove the specified policy assignment using its ID.
func Delete(c *golangsdk.ServiceClient, domainId, assignmentId string) error {
	_, err := c.Delete(resourceURL(c, domainId, assignmentId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// ListDefinitions is a method to query the list of the built-in policy definitions using given parameters.
func ListDefinitions(client *golangsdk.ServiceClient) ([]PolicyDefinition, error) {
	pages, err := pagination.NewPager(client, queryDefinitionURL(client), func(r pagination.PageResult) pagination.Page {
		p := DefinitionPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return extractDefinitions(pages)
}
