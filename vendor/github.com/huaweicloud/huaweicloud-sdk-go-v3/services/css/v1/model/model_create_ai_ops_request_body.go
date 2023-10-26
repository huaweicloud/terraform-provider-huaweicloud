package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateAiOpsRequestBody struct {

	// 检测任务名称。
	Name string `json:"name"`

	// 检测任务描述。
	Description *string `json:"description,omitempty"`

	Alarm *CreateAiOpsRequestBodyAlarm `json:"alarm,omitempty"`
}

func (o CreateAiOpsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAiOpsRequestBody struct{}"
	}

	return strings.Join([]string{"CreateAiOpsRequestBody", string(data)}, " ")
}
