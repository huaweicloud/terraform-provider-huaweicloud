package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteMessageDiagnosisReportsResponse Response Object
type BatchDeleteMessageDiagnosisReportsResponse struct {

	// 诊断报告删除结果
	Results        *[]BatchDeleteMessageDiagnosisRespResults `json:"results,omitempty"`
	HttpStatusCode int                                       `json:"-"`
}

func (o BatchDeleteMessageDiagnosisReportsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteMessageDiagnosisReportsResponse struct{}"
	}

	return strings.Join([]string{"BatchDeleteMessageDiagnosisReportsResponse", string(data)}, " ")
}
