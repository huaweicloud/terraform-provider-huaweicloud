package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DebugCaseReturnHeader struct {

	// 连接
	Connection *string `json:"Connection,omitempty"`

	// 内容长度
	ContentLength *string `json:"Content-Length,omitempty"`

	// 内容类型
	ContentType *string `json:"Content-Type,omitempty"`

	// 时间
	Date *string `json:"Date,omitempty"`

	// 兼容性保留，当前版本未使用
	Vary *string `json:"Vary,omitempty"`
}

func (o DebugCaseReturnHeader) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DebugCaseReturnHeader struct{}"
	}

	return strings.Join([]string{"DebugCaseReturnHeader", string(data)}, " ")
}
