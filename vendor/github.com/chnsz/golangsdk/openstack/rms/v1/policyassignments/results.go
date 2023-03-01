package policyassignments

import "github.com/chnsz/golangsdk/pagination"

// Assignment is the structure that represents the detail of the policy assignment.
type Assignment struct {
	// The type of the policy assignment.
	Type string `json:"policy_assignment_type"`
	// The ID of the policy assignment.
	ID string `json:"id"`
	// The name of the policy assignment.
	Name string `json:"name"`
	// The description of the policy assignment.
	Description string `json:"description"`
	// The configuration used to filter resources.
	PolicyFilter PolicyFilter `json:"policy_filter"`
	// The period of the policy rule check.
	Period string `json:"period"`
	// The status of the policy assignment.
	Status string `json:"state"`
	// The ID of the policy definition.
	PolicyDefinitionId string `json:"policy_definition_id"`
	// The configuration of the custom policy.
	CustomPolicy CustomPolicy `json:"custom_policy"`
	// The rule definition of the policy assignment.
	Parameters map[string]PolicyParameterValue `json:"parameters"`
	// The creation time of the policy assignment.
	CreatedAt string `json:"created"`
	// The latest update time of the policy assignment.
	UpdatedAt string `json:"updated"`
}

// listDefinitionResp is the structure that represents the page details of the built-in policy definitions.
type listDefinitionResp struct {
	Definitions []PolicyDefinition `json:"value"`
	// The information of the current query page.
	PageInfo pageInfo `json:"page_info"`
}

// PolicyDefinition is the structure that represents the details of the built-in policy definition.
type PolicyDefinition struct {
	// The ID of the policy definition.
	ID string `json:"id"`
	// The name of the policy definition.
	Name string `json:"name"`
	// The policy type of the policy definition.
	PolicyType string `json:"policy_type"`
	// The description of the policy definition.
	Description string `json:"description"`
	// The rule type of the policy definition.
	PolicyRuleType string `json:"policy_rule_type"`
	// The rule details of the policy definition.
	PolicyRule interface{} `json:"policy_rule"`
	// The trigger type of the policy definition.
	TriggerType string `json:"trigger_type"`
	// The keywords that policy definition has.
	Keywords []string `json:"keywords"`
	// The parameters of the policy definition.
	Parameters map[string]PolicyParameterDefinition `json:"parameters"`
}

// PolicyParameterDefinition is the structure that represents the parameter configuration.
type PolicyParameterDefinition struct {
	// The name of the parameter definition.
	Name string `json:"name"`
	// The description of the parameter definition.
	Description string `json:"description"`
	// The allow value list of the parameter definition.
	AllowedValues []interface{} `json:"allowed_values"`
	// The default value of the parameter definition.m
	DefaultValue string `json:"default_value"`
	// The type of the parameter definition.
	Type string `json:"type"`
}

// pageInfo is the structure that represents the information of the policy definition page.
type pageInfo struct {
	// The policy definition count of the current page.
	CurrentCount int `json:"current_count"`
	// The next marker of the policy definition page.
	NextMarker string `json:"next_marker"`
}

// DefinitionPage represents the response pages of the List method.
type DefinitionPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a query result no policy definition.
func (r DefinitionPage) IsEmpty() (bool, error) {
	resp, err := extractDefinitions(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index in a query result.
func (r DefinitionPage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}
	if resp.NextMarker != "" {
		return "", nil
	}
	return resp.NextMarker, nil
}

// extractPageInfo is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (*pageInfo, error) {
	var s listDefinitionResp
	err := r.(DefinitionPage).Result.ExtractInto(&s)
	return &s.PageInfo, err
}

// extractDefinitions is a method which to extract the response to a policy definition list.
func extractDefinitions(r pagination.Page) ([]PolicyDefinition, error) {
	var s listDefinitionResp
	err := r.(DefinitionPage).Result.ExtractInto(&s)
	return s.Definitions, err
}
