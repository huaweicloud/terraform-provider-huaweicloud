package entity

type ResourceImportDetailParam struct {
	ResourceId     string `json:"resource_id"`
	ResourceName   string `json:"resource_name"`
	ResourceRegion string `json:"resource_region"`
	ProjectId      string `json:"project_id"`
	EpsId          string `json:"eps_id,omitempty"`
	EpsName        string `json:"eps_name,omitempty"`
}

type ResourceImportParam struct {
	Resources []ResourceImportDetailParam `json:"resources"`
	EnvId     string                      `json:"env_id"`
}

type UnbindResourceParam struct {
	Id     string   `json:"id"`
	EnvIds []string `json:"env_ids"`
}

type DeleteResourceParam struct {
	Data []UnbindResourceParam `json:"data"`
}

type ResourceImportDetailVo struct {
	Id         string `json:"id,omitempty"`
	ResourceId string `json:"resource_id,omitempty"`
}

type CreateResourceResponse struct {
	ResourceDetail []ResourceImportDetailVo `json:"data,omitempty"`
	HttpStatusCode int                      `json:"-"`
}

type ReadResourceDetailVo struct {
	Id             string `json:"id,omitempty"`
	ResourceId     string `json:"resource_id,omitempty"`
	ResourceName   string `json:"resource_name,omitempty"`
	ResourceRegion string `json:"resource_region,omitempty"`
	EpsId          string `json:"eps_id,omitempty"`
	EpsName        string `json:"eps_name,omitempty"`
}

type ReadResourceResponse struct {
	ResourceDetail []ReadResourceDetailVo `json:"data,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

type PageResourceListParam struct {
	Maker           string            `json:"maker,omitempty"`
	Limit           string            `json:"limit,omitempty"`
	Keywords        map[string]string `json:"keywords,omitempty"`
	CiRelationships bool              `json:"ci_relationships,omitempty"`
	CiType          string            `json:"ci_type"`
	CiRegion        string            `json:"ci_region"`
	CiId            string            `json:"ci_id"`
}
