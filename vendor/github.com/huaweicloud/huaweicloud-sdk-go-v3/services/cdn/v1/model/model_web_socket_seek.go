package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// WebSocketSeek webSocket配置。  > 只有全站加速的域名支持该配置。
type WebSocketSeek struct {

	// 开关， on 开启，off 关闭。
	Status string `json:"status"`

	// 请求建立连接后，会话的保持时间：范围：1-300，单位：秒。
	Timeout int32 `json:"timeout"`
}

func (o WebSocketSeek) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "WebSocketSeek struct{}"
	}

	return strings.Join([]string{"WebSocketSeek", string(data)}, " ")
}
