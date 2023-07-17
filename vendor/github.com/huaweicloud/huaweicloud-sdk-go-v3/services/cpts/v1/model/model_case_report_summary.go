package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CaseReportSummary struct {

	// 用例和aw信息视图
	CaseAwInfoList *[]CaseAwInfo `json:"case_aw_info_list,omitempty"`

	// 错误信息
	ErrMessage *string `json:"err_message,omitempty"`
}

func (o CaseReportSummary) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CaseReportSummary struct{}"
	}

	return strings.Join([]string{"CaseReportSummary", string(data)}, " ")
}
