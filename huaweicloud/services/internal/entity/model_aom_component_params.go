package entity

type ComponentParam struct {
	Description string `json:"description,omitempty"`
	ModelId     string `json:"model_id"`
	ModelType   string `json:"model_type"`
	Name        string `json:"name"`
}

type ComponentVo struct {
	AomId        string `json:"aom_id,omitempty"`
	AppId        string `json:"app_id,omitempty"`
	CreateTime   string `json:"create_time,omitempty"`
	Creator      string `json:"creator,omitempty"`
	Description  string `json:"description,omitempty"`
	Id           string `json:"id,omitempty"`
	ModifiedTime string `json:"modified_time,omitempty"`
	Modifier     string `json:"modifier,omitempty"`
	Name         string `json:"name,omitempty"`
	RegisterType string `json:"register_type,omitempty"`
	SubAppId     string `json:"sub_app_id,omitempty"`
}
