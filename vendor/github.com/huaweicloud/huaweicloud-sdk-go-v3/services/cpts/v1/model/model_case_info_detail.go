package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CaseInfoDetail struct {

	// case_id
	CaseId *int32 `json:"case_id,omitempty"`

	// 用例名称
	Name *string `json:"name,omitempty"`

	// case_type
	CaseType *int32 `json:"case_type,omitempty"`

	// contents
	Contents *[]Contents `json:"contents,omitempty"`

	// for_loop_params
	ForLoopParams *[]interface{} `json:"for_loop_params,omitempty"`

	// increase_setting
	IncreaseSetting *[]interface{} `json:"increase_setting,omitempty"`

	// stages
	Stages *[]TestCaseStage `json:"stages,omitempty"`

	// 状态，0：已删除；1：启用；2：禁用
	Status *int32 `json:"status,omitempty"`

	// temp_id
	TempId *int32 `json:"temp_id,omitempty"`

	// 排序字段
	Sort *int32 `json:"sort,omitempty"`

	// 用例所属目录id（旧版接口可不传）
	DirectoryId *int32 `json:"directory_id,omitempty"`

	// 前置步骤
	SetupContents *[]Contents `json:"setup_contents,omitempty"`

	// 执行器个数
	UserReplicas *int32 `json:"user_replicas,omitempty"`

	// 日志采集策略（0-请求模式；1-用例模式）
	CollectLogPolicy *int32 `json:"collect_log_policy,omitempty"`

	// 关联全链路应用列表
	LinkAppList *[]int32 `json:"link_app_list,omitempty"`

	CaseInfo *CaseDoc `json:"case_info,omitempty"`
}

func (o CaseInfoDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CaseInfoDetail struct{}"
	}

	return strings.Join([]string{"CaseInfoDetail", string(data)}, " ")
}
