package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateAiOpsRequestBodyAlarm 检测任务完成后发送SMN告警消息。
type CreateAiOpsRequestBodyAlarm struct {

	// SMN告警消息敏感度。 - high：高风险。 - medium：中风险。 - suggestion：建议。 - norisk：无风险。
	Level string `json:"level"`

	// SMN主题名称。
	SmnTopic string `json:"smn_topic"`
}

func (o CreateAiOpsRequestBodyAlarm) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAiOpsRequestBodyAlarm struct{}"
	}

	return strings.Join([]string{"CreateAiOpsRequestBodyAlarm", string(data)}, " ")
}
