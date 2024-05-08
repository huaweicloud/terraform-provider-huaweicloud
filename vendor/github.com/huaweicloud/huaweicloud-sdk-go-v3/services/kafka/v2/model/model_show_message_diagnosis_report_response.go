package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowMessageDiagnosisReportResponse Response Object
type ShowMessageDiagnosisReportResponse struct {

	// 诊断异常的诊断项总和
	AbnormalItemNum *int32 `json:"abnormal_item_num,omitempty"`

	// 诊断失败的诊断项总和
	FailedItemNum *int32 `json:"failed_item_num,omitempty"`

	// 诊断正常的诊断项总和
	NormalItemNum *int32 `json:"normal_item_num,omitempty"`

	// 诊断维度列表
	DiagnosisDimensionList *[]KafkaMessageDiagnosisDimensionEntity `json:"diagnosis_dimension_list,omitempty"`
	HttpStatusCode         int                                     `json:"-"`
}

func (o ShowMessageDiagnosisReportResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMessageDiagnosisReportResponse struct{}"
	}

	return strings.Join([]string{"ShowMessageDiagnosisReportResponse", string(data)}, " ")
}
