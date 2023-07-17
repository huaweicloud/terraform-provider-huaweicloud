package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ReportOutlineResult struct {

	// 错误信息
	ErrMessage *interface{} `json:"err_message,omitempty"`

	Outline *ReportOutline `json:"outline,omitempty"`
}

func (o ReportOutlineResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReportOutlineResult struct{}"
	}

	return strings.Join([]string{"ReportOutlineResult", string(data)}, " ")
}
