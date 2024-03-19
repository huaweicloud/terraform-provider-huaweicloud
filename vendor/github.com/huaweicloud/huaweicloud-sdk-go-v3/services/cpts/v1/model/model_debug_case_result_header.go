package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DebugCaseResultHeader struct {

	// 连接
	Connection *string `json:"Connection,omitempty"`

	// 内容类型
	ContentType *string `json:"Content-Type,omitempty"`

	// 主机
	Host *string `json:"Host,omitempty"`
}

func (o DebugCaseResultHeader) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DebugCaseResultHeader struct{}"
	}

	return strings.Join([]string{"DebugCaseResultHeader", string(data)}, " ")
}
