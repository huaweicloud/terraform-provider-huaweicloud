package applications

import "github.com/chnsz/golangsdk/pagination"

// Application is the structure that represents the detail of the ServiceStage application.
type Application struct {
	// The application ID.
	ID string `json:"id"`
	// The application name.
	Name string `json:"name"`
	// The application description.
	Description string `json:"description"`
	// The creator.
	Creator string `json:"creator"`
	// Thec project ID.
	ProjectId string `json:"project_id"`
	// The enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// The creation time.
	CreatedAt int `json:"create_time"`
	// The update time.
	UpdatedAt int `json:"update_time"`
	// Whether to enable the unified model.
	UnifiedModel string `json:"unified_model"`
}

// ApplicationPage is a single page maximum result representing a query by offset page.
type ApplicationPage struct {
	pagination.OffsetPageBase
}

// ExtractApplications is a method to extract the list of application details for ServiceStage service.
func ExtractApplications(r pagination.Page) ([]Application, error) {
	var s []Application
	r.(ApplicationPage).Result.ExtractIntoSlicePtr(&s, "applications")
	return s, nil
}

// ConfigResp is the structure that represents all configurations in the specified environment.
type ConfigResp struct {
	// The application ID.
	ApplicationId string `json:"application_id"`
	// The environment ID.
	EnvironmentId string `json:"environment_id"`
	// The configurations of the application.
	Configuration Configuration `json:"configuration"`
}
