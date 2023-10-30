package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ConnectState 连接端状态
type ConnectState struct {

	// 隧道最近一次状态更新时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	LastUpdateTime *string `json:"last_update_time,omitempty"`

	// 客户端连接状态 CONNECTED | DISCONNECTED
	Status *string `json:"status,omitempty"`
}

func (o ConnectState) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConnectState struct{}"
	}

	return strings.Join([]string{"ConnectState", string(data)}, " ")
}
