package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PortStatisticResponseInfo 开放端口统计信息
type PortStatisticResponseInfo struct {

	// 端口号
	Port *int32 `json:"port,omitempty"`

	// 类型
	Type *string `json:"type,omitempty"`

	// 端口数量
	Num *int32 `json:"num,omitempty"`
}

func (o PortStatisticResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PortStatisticResponseInfo struct{}"
	}

	return strings.Join([]string{"PortStatisticResponseInfo", string(data)}, " ")
}
