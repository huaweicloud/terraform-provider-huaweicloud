package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDiagnosisPreCheckResponse Response Object
type ShowDiagnosisPreCheckResponse struct {

	// Kafka消息积压诊断预检查返回对象
	Body           *[]KafkaDiagnosisCheckEntity `json:"body,omitempty"`
	HttpStatusCode int                          `json:"-"`
}

func (o ShowDiagnosisPreCheckResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDiagnosisPreCheckResponse struct{}"
	}

	return strings.Join([]string{"ShowDiagnosisPreCheckResponse", string(data)}, " ")
}
