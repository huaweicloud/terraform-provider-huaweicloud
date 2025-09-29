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
	// The environment type.
	DeployMode string `json:"deploy_mode"`
	// The environment description.
	Description string `json:"description"`
	// The project ID.
	ProjectId string `json:"project_id"`
	// The enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// The charging mode.
	ChargeMode string `json:"charge_mode"`
	// The VPC ID.
	VpcId string `json:"vpc_id"`
	// The basic resources.
	BaseResources []Resource `json:"base_resources"`
	// The optional resources.
	OptionalResources []Resource `json:"optional_resources"`
	// The Creator.
	Creator string `json:"creator"`
	// The creation time.
	CreatedAt int `json:"create_time"`
	// The update time.
	UpdatedAt int `json:"update_time"`
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
