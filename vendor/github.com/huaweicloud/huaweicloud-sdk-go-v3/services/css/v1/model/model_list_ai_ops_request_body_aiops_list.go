package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ListAiOpsRequestBodyAiopsList struct {

	// 检测任务id。
	Id *string `json:"id,omitempty"`

	// 检测任务名称。
	Name *string `json:"name,omitempty"`

	// 检测任务描述。
	Desc *string `json:"desc,omitempty"`

	// 任务执行状态。 - 150：未开启。 - 200：已开启。 - 300：已发送。
	Status *int32 `json:"status,omitempty"`

	Summary *ListAiOpsRequestBodySummary `json:"summary,omitempty"`

	// 检测任务创建时间戳。
	CreateTime *string `json:"create_time,omitempty"`

	// 检测任务SMN告警任务发送状态。 - not_open：未开启。 - not_trigger：未触发。 - sent：已发送。 - send_fail： 发送失败。
	SmnStatus *string `json:"smn_status,omitempty"`

	// 发送失败原因。
	SmnFailReason *string `json:"smn_fail_reason,omitempty"`

	// 风险项详情。
	TaskRisks *[]AiOpsRiskInfo `json:"task_risks,omitempty"`
}

func (o ListAiOpsRequestBodyAiopsList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAiOpsRequestBodyAiopsList struct{}"
	}

	return strings.Join([]string{"ListAiOpsRequestBodyAiopsList", string(data)}, " ")
}
