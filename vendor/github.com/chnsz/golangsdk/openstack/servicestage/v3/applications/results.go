package applications

import "github.com/chnsz/golangsdk/pagination"

// Application is the structure that represents the detail of the ServiceStage application.
type Application struct {
	// The application ID.
	ID string `json:"id"`
	// The application name.
	Name string `json:"name"`
	// The number of components under application
	ComponentCount int `json:"component_count"`
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
	// The application description.
	Description string `json:"description"`
	// Whether to enable the unified model.
	Labels []string `json:"labels"`
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
