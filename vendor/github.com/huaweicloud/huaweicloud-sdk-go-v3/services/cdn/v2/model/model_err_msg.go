package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ErrMsg 用于返回具体的错误码和错误消息
type ErrMsg struct {

	// 错误码
	ErrorCode string `json:"error_code"`

	// 错误描述
	ErrorMsg string `json:"error_msg"`
}

func (o ErrMsg) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ErrMsg struct{}"
	}

	return strings.Join([]string{"ErrMsg", string(data)}, " ")
}
