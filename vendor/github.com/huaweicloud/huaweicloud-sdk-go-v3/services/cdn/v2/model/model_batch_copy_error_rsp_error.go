package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCopyErrorRspError 错误体
type BatchCopyErrorRspError struct {

	// 错误码
	ErrorCode *string `json:"error_code,omitempty"`

	// 错误描述
	ErrorMsg *string `json:"error_msg,omitempty"`
}

func (o BatchCopyErrorRspError) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCopyErrorRspError struct{}"
	}

	return strings.Join([]string{"BatchCopyErrorRspError", string(data)}, " ")
}
