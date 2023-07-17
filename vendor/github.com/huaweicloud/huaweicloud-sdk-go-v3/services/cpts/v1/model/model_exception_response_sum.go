package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ExceptionResponseSum struct {

	// 断言失败数
	FailedAssert *int32 `json:"failed_assert,omitempty"`

	// 其他错误失败数
	FailedOthers *int32 `json:"failed_others,omitempty"`

	// 解析失败数
	FailedParsed *int32 `json:"failed_parsed,omitempty"`

	// 连接被拒绝失败数
	FailedRefused *int32 `json:"failed_refused,omitempty"`

	// 响应超时失败数
	FailedTimeout *int32 `json:"failed_timeout,omitempty"`
}

func (o ExceptionResponseSum) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExceptionResponseSum struct{}"
	}

	return strings.Join([]string{"ExceptionResponseSum", string(data)}, " ")
}
