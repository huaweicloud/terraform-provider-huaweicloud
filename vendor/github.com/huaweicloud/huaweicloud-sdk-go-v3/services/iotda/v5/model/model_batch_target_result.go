package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchTargetResult 批量操作目标结果
type BatchTargetResult struct {

	// 执行批量任务的目标。
	Target *string `json:"target,omitempty"`

	// 目标的执行结果，为success或failure
	Status *string `json:"status,omitempty"`

	// 操作失败的错误码
	ErrorCode *string `json:"error_code,omitempty"`

	// 操作失败的错误描述
	ErrorMsg *string `json:"error_msg,omitempty"`
}

func (o BatchTargetResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchTargetResult struct{}"
	}

	return strings.Join([]string{"BatchTargetResult", string(data)}, " ")
}
