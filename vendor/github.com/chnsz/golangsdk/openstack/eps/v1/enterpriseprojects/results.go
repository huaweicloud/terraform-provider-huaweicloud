package enterpriseprojects

import "github.com/chnsz/golangsdk"

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Type        string `json:"type"`
}

type Projects struct {
	EnterpriseProjects []Project `json:"enterprise_projects"`
	TotalCount         int       `json:"total_count"`
}

type ListResult struct {
	golangsdk.Result
}

func (r ListResult) Extract() ([]Project, error) {
	var a struct {
		EnterpriseProjects []Project `json:"enterprise_projects"`
	}
	err := r.Result.ExtractInto(&a)
	return a.EnterpriseProjects, err
}

type commonResult struct {
	golangsdk.Result
}

type CreatResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type ActionResult struct {
	commonResult
}

type MigrateResult struct {
	commonResult
}

type ResourceResult struct {
	commonResult
}

type FilterResult struct {
	Resources  []Resource `json:"resources"`
	Errors     []Errors   `json:"errors"`
	TotalCount int32      `json:"total_count"`
}

type Resource struct {
	EnterpriseProjectId string `json:"enterprise_project_id"`

	ProjectId string `json:"project_id"`

	ProjectName string `json:"project_name"`

	ResourceDetail interface{} `json:"-"`

	ResourceId string `json:"resource_id"`

	ResourceName string `json:"resource_name"`

	ResourceType string `json:"resource_type"`
}

type Errors struct {
	ErrorCode string `json:"error_code,omitempty"`

	ErrorMsg string `json:"error_msg,omitempty"`

	ProjectId string `json:"project_id,omitempty"`

	ResourceType string `json:"resource_type,omitempty"`
}

func (r commonResult) Extract() (Project, error) {
	var s struct {
		EnterpriseProject Project `json:"enterprise_project"`
	}
	err := r.ExtractInto(&s)
	return s.EnterpriseProject, err
}
