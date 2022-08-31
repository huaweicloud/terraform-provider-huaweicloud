package entity

type DashBoardRequest struct {
	// 日志组id
	LogGroupId string `json:"log_group_id"`

	// 目标日志组名称。
	LogGroupName string `json:"log_group_name"`

	// 日志流id
	LogStreamId string `json:"log_stream_id"`
	// 目标日志组名称。
	LogStreamName string `json:"log_stream_name"`

	TemplateTitle []string `json:"template_title"`
	TemplateType  []string `json:"template_type"`
	GroupName     string   `json:"group_name"`
}

type ReadDashBoardResp struct {
	Results []DashBoard `json:"results"`
}

type DashBoard struct {
	ProjectId         string           `json:"project_id"`
	Id                string           `json:"id"`
	GroupName         string           `json:"group_name"`
	Title             string           `json:"title"`
	Charts            []DashboardChars `json:"charts"`
	Filters           []interface{}    `json:"filters"`
	LastUpdateTime    int              `json:"last_update_time"`
	UseSystemTemplate bool             `json:"useSystemTemplate"`
}

type DashboardChars struct {
	Width   int                    `json:"width"`
	Height  int                    `json:"height"`
	XPos    int                    `json:"x_pos"`
	YPos    int                    `json:"y_pos"`
	ChartId string                 `json:"chart_id"`
	Chart   map[string]interface{} `json:"chart"`
}
