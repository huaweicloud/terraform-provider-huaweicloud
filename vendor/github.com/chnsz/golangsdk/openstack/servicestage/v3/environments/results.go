package environments

import "github.com/chnsz/golangsdk/pagination"

// Environment is the structure that represents the detail of the ServiceStage environment.
type Environment struct {
	// The environment ID.
	ID string `json:"id"`
	// The environment name.
	Name string `json:"name"`
	// The environment alias
	Alias string `json:"alias"`
	// The environment description.
	Description string `json:"description"`
	// The project ID.
	ProjectId string `json:"project_id"`
	// The enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// Specified the billing mode. The valid values are:
	//   provided: provided resources are used and no fees are charged.
	//   on_demanded: on-demand charging.
	//   monthly: monthly subscription.
	ChargeMode string `json:"charge_mode,omitempty"`
	// Environment deploy mode. Value: container, virtualmachine or mixed.
	DeployMode string `json:"deploy_mode"`
	// The VPC ID.
	VpcId string `json:"vpc_id"`
	// The resources info in the environment.
	Resources []Resource `json:"resources"`
	// The Creator.
	Creator string `json:"creator"`
	// The creation time.
	CreatedAt int `json:"create_time"`
	// The update time.
	UpdatedAt int `json:"update_time"`
	// It only takes effect when the environment's DeployMode is "virtualmachine"
	// Value: 50 or 500
	VmClusterSize int `json:"vm_cluster_size"`
	// Environment's Labels
	Labels []Label `json:"labels"`
	// It only takes effect when environment's DeployMode is "virtualmachine" And VmClusterSize is 500
	BrokerID string `json:"broker_id"`
	// Number of components in the environment
	ComponentCount int `json:"component_count"`
}

// Label is the structure that represents some label of environment.
type Label struct {
	// Label Key
	Key string `json:"key"`
	// Label Value
	Value string `json:"value"`
}

// EnvironmentPage is a single page maximum result representing a query by offset page.
type EnvironmentPage struct {
	pagination.OffsetPageBase
}

// ExtractEnvironments is a method to extract the list of environment details for ServiceStage service.
func ExtractEnvironments(r pagination.Page) ([]Environment, error) {
	var s []Environment
	r.(EnvironmentPage).Result.ExtractIntoSlicePtr(&s, "environments")
	return s, nil
}

// Resource is an object specifying the basic or optional resource.
type Resource struct {
	// Specified the resource ID.
	ID string `json:"id" required:"true"`
	// Specified the resource type. the valid values are: CCE, CCI, ECS and AS.
	Type string `json:"type" required:"true"`
	// Specified the resource name.
	Name string `json:"name,omitempty"`
}
