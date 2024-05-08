package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// KafkaMessageDiagnosisDimensionEntity 消息积压诊断维度
type KafkaMessageDiagnosisDimensionEntity struct {

	// 诊断维度名称
	Name string `json:"name"`

	// 该诊断维度下，异常的诊断项总数
	AbnormalNum int32 `json:"abnormal_num"`

	// 该诊断维度下，诊断失败的诊断项总和
	FailedNum int32 `json:"failed_num"`

	// 诊断项列表
	DiagnosisItemList []KafkaMessageDiagnosisItemEntity `json:"diagnosis_item_list"`
}

func (o KafkaMessageDiagnosisDimensionEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KafkaMessageDiagnosisDimensionEntity struct{}"
	}

	return strings.Join([]string{"KafkaMessageDiagnosisDimensionEntity", string(data)}, " ")
}
