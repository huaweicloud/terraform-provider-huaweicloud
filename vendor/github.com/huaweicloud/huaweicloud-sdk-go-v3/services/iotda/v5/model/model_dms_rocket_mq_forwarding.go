package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DmsRocketMqForwarding rocketMQ服务配置信息
type DmsRocketMqForwarding struct {

	// **参数说明**：转发rocketMQ消息对应的地址列表
	Addresses []NetAddress `json:"addresses"`

	// **参数说明**：转发rocketMQ消息关联的topic信息。
	Topic string `json:"topic"`

	// **参数说明**：转发rocketMQ关联的用户名信息。
	Username string `json:"username"`

	// **参数说明**：转发rocketMQ关联的密码信息。
	Password string `json:"password"`

	// 是否开启SSL，默认为true。
	EnableSsl *bool `json:"enable_ssl,omitempty"`
}

func (o DmsRocketMqForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DmsRocketMqForwarding struct{}"
	}

	return strings.Join([]string{"DmsRocketMqForwarding", string(data)}, " ")
}
