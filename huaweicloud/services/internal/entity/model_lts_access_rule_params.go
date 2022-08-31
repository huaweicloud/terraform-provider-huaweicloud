package entity

type AomMappingRequestInfo struct {

	// 项目id
	ProjectId string `json:"project_id"`

	// 接入规则名称
	RuleName string `json:"rule_name"`

	// 接入规则id
	RuleId string `json:"rule_id,omitempty"`

	RuleInfo AomMappingRuleInfo `json:"rule_info"`
}

type AomMappingRuleInfo struct {

	// 集群id
	ClusterId string `json:"cluster_id"`

	// 集群名称
	ClusterName string `json:"cluster_name"`

	// 日志流前缀
	DeploymentsPrefix *string `json:"deployments_prefix,omitempty"`

	// 工作负载
	Deployments []string `json:"deployments"`

	// 命名空间
	Namespace string `json:"namespace"`

	// 容器名称
	ContainerName *string `json:"container_name,omitempty"`

	// 接入规则详情
	Files []AomMappingfilesInfo `json:"files"`
}

type AomMappingfilesInfo struct {

	// 路径名
	FileName string `json:"file_name"`

	LogStreamInfo AomMappingLogStreamInfo `json:"log_stream_info"`
}

type AomMappingLogStreamInfo struct {

	// 日志组id
	TargetLogGroupId string `json:"target_log_group_id"`

	// 目标日志组名称。
	TargetLogGroupName string `json:"target_log_group_name"`

	// 日志流id
	TargetLogStreamId string `json:"target_log_stream_id"`

	// 目标日志组名称。
	TargetLogStreamName string `json:"target_log_stream_name"`
}

type CreateAomMappingRulesResponse struct {
	Body           *[]AomMappingRuleResp `json:"body,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

type AomMappingRuleResp struct {

	// 项目id
	ProjectId string `json:"project_id"`

	// 接入规则名称
	RuleName string `json:"rule_name"`

	// 接入规则id
	RuleId string `json:"rule_id"`

	RuleInfo *AomMappingRuleInfo `json:"rule_info"`
}
