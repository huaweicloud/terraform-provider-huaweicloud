package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DevicePolicyBindOrUnbindFailureDetail 策略绑定失败详情结构体。
type DevicePolicyBindOrUnbindFailureDetail struct {

	// 失败的目标id。
	TargetId *string `json:"target_id,omitempty"`

	// 错误码。
	ErrorCode *string `json:"error_code,omitempty"`

	// 错误详情。
	ErrorMsg *string `json:"error_msg,omitempty"`
}

func (o DevicePolicyBindOrUnbindFailureDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DevicePolicyBindOrUnbindFailureDetail struct{}"
	}

	return strings.Join([]string{"DevicePolicyBindOrUnbindFailureDetail", string(data)}, " ")
}
