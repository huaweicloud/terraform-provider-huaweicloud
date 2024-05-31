package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResetBridgeSecret struct {

	// 是否强制断开网桥的连接，当前仅限长连接。
	ForceDisconnect *bool `json:"force_disconnect,omitempty"`
}

func (o ResetBridgeSecret) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetBridgeSecret struct{}"
	}

	return strings.Join([]string{"ResetBridgeSecret", string(data)}, " ")
}
