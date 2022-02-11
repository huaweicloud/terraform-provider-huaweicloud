package environments

import "github.com/chnsz/golangsdk/pagination"

type ListResp struct {
	// Number of environments that match the query conditions.
	Total int `json:"total"`
	// Length of the returned environment list.
	Size int `json:"size"`
	// Environment list.
	Environments []Environment `json:"envs"`
}

type Environment struct {
	// Environment ID.
	Id string `json:"id"`
	// Environment name.
	Name string `json:"name"`
	// Time when the environment is created.
	CreateTime string `json:"create_time"`
	// Description of the environment.
	Description string `json:"remark"`
}

// EnvironmentPage represents the response pages of the List method.
type EnvironmentPage struct {
	pagination.SinglePageBase
}
