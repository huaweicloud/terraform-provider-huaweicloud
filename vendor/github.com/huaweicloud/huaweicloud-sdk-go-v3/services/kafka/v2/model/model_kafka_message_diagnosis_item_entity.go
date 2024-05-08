package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// KafkaMessageDiagnosisItemEntity 消息积压诊断项
type KafkaMessageDiagnosisItemEntity struct {

	// 诊断项名称
	Name string `json:"name"`

	// 诊断结果
	Result string `json:"result"`

	// 诊断异常原因列表
	CauseIds *[]KafkaMessageDiagnosisConclusionEntity `json:"cause_ids,omitempty"`

	// 诊断异常建议列表
	AdviceIds *[]KafkaMessageDiagnosisConclusionEntity `json:"advice_ids,omitempty"`

	// 诊断异常受影响的分区列表
	Partitions *[]int32 `json:"partitions,omitempty"`

	// 诊断失败的分区列表
	FailedPartitions *[]int32 `json:"failed_partitions,omitempty"`

	// 诊断异常受影响的broker列表
	BrokerIds *[]int32 `json:"broker_ids,omitempty"`
}

func (o KafkaMessageDiagnosisItemEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KafkaMessageDiagnosisItemEntity struct{}"
	}

	return strings.Join([]string{"KafkaMessageDiagnosisItemEntity", string(data)}, " ")
}
