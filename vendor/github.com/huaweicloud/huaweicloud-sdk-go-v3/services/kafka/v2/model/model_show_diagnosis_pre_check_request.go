package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDiagnosisPreCheckRequest Request Object
type ShowDiagnosisPreCheckRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 消费组名称
	Group string `json:"group"`

	// 主题名称
	Topic string `json:"topic"`
}

func (o ShowDiagnosisPreCheckRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDiagnosisPreCheckRequest struct{}"
	}

	return strings.Join([]string{"ShowDiagnosisPreCheckRequest", string(data)}, " ")
}
