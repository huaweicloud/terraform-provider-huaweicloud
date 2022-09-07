package entity

type EnvParam struct {
	ComponentId  string `json:"component_id,omitempty"`
	Description  string `json:"description,omitempty"`
	EnvName      string `json:"env_name"`
	EnvType      string `json:"env_type"`
	OsType       string `json:"os_type,omitempty"`
	Region       string `json:"region"`
	RegisterType string `json:"register_type,omitempty"`
}

type EnvVo struct {
	AomId        string           `json:"aom_id,omitempty"`
	ComponentId  string           `json:"component_id,omitempty"`
	CreateTime   string           `json:"create_time,omitempty"`
	Creator      string           `json:"creator,omitempty"`
	Description  string           `json:"description,omitempty"`
	EnvId        string           `json:"env_id,omitempty"`
	EnvName      string           `json:"env_name,omitempty"`
	EnvTags      []TagNameAndIdVo `json:"env_tags,omitempty"`
	EnvType      string           `json:"env_type,omitempty"`
	EpsId        string           `json:"eps_id,omitempty"`
	ModifiedTime string           `json:"modified_time,omitempty"`
	Modifier     string           `json:"modifier,omitempty"`
	OsType       string           `json:"os_type,omitempty"`
	Region       string           `json:"region,omitempty"`
	RegisterType string           `json:"register_type,omitempty"`
}

type TagNameAndIdVo struct {
	TagId   string `json:"tag_id,omitempty"`
	TagName string `json:"tag_name,omitempty"`
}
