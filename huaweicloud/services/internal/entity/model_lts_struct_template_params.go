package entity

type StructTemplateRequest struct {
	Content      string           `json:"content,omitempty"`
	LogGroupId   string           `json:"log_group_id"`
	ParseType    string           `json:"parse_type,omitempty"`
	TemplateId   string           `json:"template_id"`
	TemplateType string           `json:"template_type"`
	TemplateName string           `json:"template_name"`
	LogStreamId  string           `json:"log_stream_id"`
	ProjectId    string           `json:"project_id"`
	RegexRules   *string          `json:"regex_rules,omitempty"`
	Layers       *int             `json:"layers,omitempty"`
	Tokenizer    string           `json:"tokenizer,omitempty"`
	LogFormat    *string          `json:"log_format,omitempty"`
	DemoFields   []DemoFieldsInfo `json:"demo_fields"`
	TagFields    []TagFieldsInfo  `json:"tag_fields"`
}

type DemoFieldsInfo struct {
	IsAnalysis      bool   `json:"is_analysis"`
	Content         string `json:"content,omitempty"`
	FieldName       string `json:"field_name,omitempty"`
	Type            string `json:"type,omitempty"`
	UserDefinedName string `json:"userDefinedName,omitempty"`
	Index           int    `json:"index,omitempty"`
}

type TagFieldsInfo struct {
	FieldName  string  `json:"fieldName"`
	Type       string  `json:"type"`
	Content    *string `json:"content,omitempty"`
	IsAnalysis *bool   `json:"isAnalysis,omitempty"`
}

func (s *StructTemplateRequest) ToDemoFieldsInfo() {
	s.DemoFields = []DemoFieldsInfo{
		{
			Type:      "string",
			FieldName: "remote_ip",
			Index:     0,
		},
		{
			Type:      "string",
			FieldName: "local_ip",
			Index:     1,
		},
		{
			Type:      "string",
			FieldName: "local_port",
			Index:     2,
		},
		{
			Type:      "string",
			FieldName: "t",
			Index:     3,
		},
		{
			Type:      "string",
			FieldName: "tt",
			Index:     4,
		},
		{
			Type:      "string",
			FieldName: "method",
			Index:     5,
		},
		{
			Type:      "string",
			FieldName: "uri",
			Index:     6,
		},
		{
			Type:      "string",
			FieldName: "protocal",
			Index:     7,
		},
		{
			Type:      "string",
			FieldName: "code",
			Index:     8,
		},
		{
			Type:      "string",
			FieldName: "send_bytes",
			Index:     9,
		},
		{
			Type:      "string",
			FieldName: "cost",
			Index:     10,
		},
	}
}

type ShowStructTemplateResponse struct {

	// 结构化字段
	DemoFields *[]StructFieldInfoReturn `json:"demoFields,omitempty"`

	// 关键词详细信息
	TagFields *[]StructTagFieldsInfo `json:"tagFields,omitempty"`

	// 示例日志
	DemoLog *string `json:"demoLog,omitempty"`

	// 测试
	DemoLabel *string `json:"demoLabel,omitempty"`

	// id
	Id string `json:"id,omitempty"`

	// 日志组ID
	LogGroupId *string `json:"logGroupId,omitempty"`

	Rule *ShowStructTemplateRule `json:"rule,omitempty"`

	// 日志流ID
	LogStreamId *string `json:"logStreamId,omitempty"`

	// 项目ID
	ProjectId *string `json:"projectId,omitempty"`

	// 测试
	TemplateName *string `json:"templateName,omitempty"`

	// 为了兼容前台数据格式
	Regex          *string `json:"regex,omitempty"`
	HttpStatusCode int     `json:"-"`
}

type StructFieldInfoReturn struct {

	// 字段名称
	FieldName *string `json:"fieldName,omitempty"`

	// 字段数据类型
	Type *string `json:"type,omitempty"`

	// 字段内容
	Content *string `json:"content,omitempty"`

	// 结构化方式
	IsAnalysis *bool `json:"isAnalysis,omitempty"`

	// 序号
	Index *int32 `json:"index,omitempty"`
}

type StructTagFieldsInfo struct {

	// 字段名称
	FieldName *string `json:"fieldName,omitempty"`

	// 字段类型
	Type *string `json:"type,omitempty"`

	// 内容
	Content *string `json:"content,omitempty"`

	// 是否解析
	IsAnalysis *bool `json:"isAnalysis,omitempty"`

	// 字段名称
	Index *int32 `json:"index,omitempty"`
}

type ShowStructTemplateRule struct {

	// 测试
	Param *string `json:"param,omitempty"`

	// 结构化类型
	Type *string `json:"type,omitempty"`
}

type DeleteStructTemplateReqBody struct {

	// 结构化规则ID
	Id string `json:"id"`
}
