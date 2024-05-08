package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// KafkaDiagnosisCheckEntity kafka消息积压诊断预检查实体
type KafkaDiagnosisCheckEntity struct {

	// 预检查项名称
	Name string `json:"name"`

	// 预检查失败原因
	Reason string `json:"reason"`

	// 预检查是否正常
	Success bool `json:"success"`
}

func (o KafkaDiagnosisCheckEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KafkaDiagnosisCheckEntity struct{}"
	}

	return strings.Join([]string{"KafkaDiagnosisCheckEntity", string(data)}, " ")
}
