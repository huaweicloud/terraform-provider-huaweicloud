package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RomaForwarding 转发roma消息结构
type RomaForwarding struct {

	// **参数说明**：转发roma消息对应的地址列表
	Addresses []NetAddress `json:"addresses"`

	// **参数说明**：转发roma消息关联的topic信息。
	Topic string `json:"topic"`

	// **参数说明**：转发roma关联的用户名信息。
	Username string `json:"username"`

	// **参数说明**：转发roma关联的密码信息。
	Password string `json:"password"`
}

func (o RomaForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RomaForwarding struct{}"
	}

	return strings.Join([]string{"RomaForwarding", string(data)}, " ")
}
