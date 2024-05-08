package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// KafkaMessageDiagnosisConclusionEntity 消息积压诊断结论
type KafkaMessageDiagnosisConclusionEntity struct {

	// 诊断结论ID
	Id int32 `json:"id"`

	// 诊断结论参数列表
	Params map[string]string `json:"params,omitempty"`
}

func (o KafkaMessageDiagnosisConclusionEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KafkaMessageDiagnosisConclusionEntity struct{}"
	}

	return strings.Join([]string{"KafkaMessageDiagnosisConclusionEntity", string(data)}, " ")
}
