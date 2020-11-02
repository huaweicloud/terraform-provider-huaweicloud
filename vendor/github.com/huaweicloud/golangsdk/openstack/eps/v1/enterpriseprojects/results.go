package enterpriseprojects

import "github.com/huaweicloud/golangsdk"

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
