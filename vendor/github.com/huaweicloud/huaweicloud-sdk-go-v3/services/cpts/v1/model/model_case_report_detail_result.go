package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CaseReportDetailResult struct {
	Detail *CaseReportDetails `json:"detail,omitempty"`

	// 错误信息
	ErrMessage *string `json:"err_message,omitempty"`
}

func (o CaseReportDetailResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CaseReportDetailResult struct{}"
	}

	return strings.Join([]string{"CaseReportDetailResult", string(data)}, " ")
}
