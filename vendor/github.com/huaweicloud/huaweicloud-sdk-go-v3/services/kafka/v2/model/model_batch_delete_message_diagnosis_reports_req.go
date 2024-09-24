package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteMessageDiagnosisReportsReq 批量删除消息积压诊断报告请求
type BatchDeleteMessageDiagnosisReportsReq struct {

	// 待删除report id列表
	ReportIdList []string `json:"report_id_list"`
}

func (o BatchDeleteMessageDiagnosisReportsReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteMessageDiagnosisReportsReq struct{}"
	}

	return strings.Join([]string{"BatchDeleteMessageDiagnosisReportsReq", string(data)}, " ")
}
