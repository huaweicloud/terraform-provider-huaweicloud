package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchDeleteMessageDiagnosisRespResults struct {

	// 报告删除结果
	Result bool `json:"result"`

	// 报告ID
	Id string `json:"id"`
}

func (o BatchDeleteMessageDiagnosisRespResults) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteMessageDiagnosisRespResults struct{}"
	}

	return strings.Join([]string{"BatchDeleteMessageDiagnosisRespResults", string(data)}, " ")
}
