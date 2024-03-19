package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CaseInfoDetail struct {

	// 用例id
	CaseId *int32 `json:"case_id,omitempty"`

	// 用例名称
	Name *string `json:"name,omitempty"`

	// 用例类型（0：常规用例；）
	CaseType *int32 `json:"case_type,omitempty"`

	// 用例脚本信息
	Contents *[]Contents `json:"contents,omitempty"`

	// 旧版本逻辑控制器字段，当前已未使用
	ForLoopParams *[]interface{} `json:"for_loop_params,omitempty"`

	// 梯度递增
	IncreaseSetting *[]interface{} `json:"increase_setting,omitempty"`

	// 阶段信息
	Stages *[]TestCaseStage `json:"stages,omitempty"`

	// 状态，0：已删除；1：启用；2：禁用
	Status *int32 `json:"status,omitempty"`

	// 用例id
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
