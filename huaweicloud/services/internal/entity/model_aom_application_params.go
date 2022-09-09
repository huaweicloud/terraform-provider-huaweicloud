package entity

type BizAppParam struct {
	Description  string `json:"description,omitempty"`
	DisplayName  string `json:"display_name,omitempty"`
	EpsId        string `json:"eps_id,omitempty"`
	Name         string `json:"name,omitempty"`
	RegisterType string `json:"register_type,omitempty"`
}

type CreateModelVo struct {
	Id string `json:"id"`
}

type BizAppVo struct {
	AomId        string `json:"aom_id,omitempty"`
	AppId        string `json:"app_id,omitempty"`
	CreateTime   string `json:"create_time,omitempty"`
	Creator      string `json:"creator,omitempty"`
	Description  string `json:"description,omitempty"`
	DisplayName  string `json:"display_name,omitempty"`
	EpsId        string `json:"eps_id,omitempty"`
	ModifiedTime string `json:"modified_time,omitempty"`
	Modifier     string `json:"modifier,omitempty"`
	Name         string `json:"name,omitempty"`
	RegisterType string `json:"register_type,omitempty"`
}
