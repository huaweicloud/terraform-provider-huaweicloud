package components

import "github.com/chnsz/golangsdk/pagination"

// Component is the structure that represents the detail of the application component.
type Component struct {
	// Application component ID.
	ID string `json:"id"`
	// Application component name.
	Name string `json:"name"`
	// Value: 0 or 1.
	// 0: Normal.
	// 1: Being deleted.
	Status int `json:"status"`
	// Runtime.
	Runtime string `json:"runtime"`
	// Application component type. Example: Webapp, MicroService, or Common.
	Type string `json:"category"`
	// Application component sub-type.
	// Webapp sub-types include Web.
	// MicroService sub-types include Java Chassis, Go Chassis, Mesher, Spring Cloud, and Dubbo.
	// Common sub-type can be empty.
	Framwork string `json:"sub_category"`
	// Description.
	Description string `json:"description"`
	// Project ID.
	ProjectId string `json:"project_id"`
	// Application ID.
	ApplicationId string `json:"application_id"`
	// Source of the code or software package.
	Source Source `json:"source"`
	// Component Builder.
	Builder Builder `json:"build"`
	// Pipeline ID list. A maximum of 10 pipeline IDs are supported.
	PipelineIds []string `json:"pipeline_ids"`
	// Creation time.
	CreatedAt int `json:"create_time"`
	// Update time.
	UpdatedAt int `json:"update_time"`
	// Creator.
	Creator string `json:"creator"`
}

// ComponentPage is a single page maximum result representing a query by offset page.
type ComponentPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a ComponentPage struct is empty.
func (b ComponentPage) IsEmpty() (bool, error) {
	arr, err := ExtractComponents(b)
	return len(arr) == 0, err
}

// ExtractComponents is a method to extract the list of component details for ServiceStage service.
func ExtractComponents(r pagination.Page) ([]Component, error) {
	var s []Component
	err := r.(ComponentPage).Result.ExtractIntoSlicePtr(&s, "components")
	return s, err
}
