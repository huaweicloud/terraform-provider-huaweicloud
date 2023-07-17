package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// MrsKafkaForwarding 转发MRS Kafka消息结构
type MrsKafkaForwarding struct {

	// **参数说明**：转发kafka消息对应的地址列表
	Addresses []NetAddress `json:"addresses"`

	// **参数说明**：转发kafka消息关联的topic信息。
	Topic string `json:"topic"`

	// 是否Kerberos认证，默认为false。
	KerberosAuthentication *bool `json:"kerberos_authentication,omitempty"`
}

func (o MrsKafkaForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MrsKafkaForwarding struct{}"
	}

	return strings.Join([]string{"MrsKafkaForwarding", string(data)}, " ")
}
