package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 转发kafka消息结构
type DmsKafkaForwarding struct {

	// **参数说明**：Kafka服务对应的region区域
	RegionName string `json:"region_name"`

	// **参数说明**：Kafka服务对应的projectId信息
	ProjectId string `json:"project_id"`

	// **参数说明**：转发kafka消息对应的地址列表
	Addresses []NetAddress `json:"addresses"`

	// **参数说明**：转发kafka消息关联的topic信息。
	Topic string `json:"topic"`

	// **参数说明**：转发kafka关联的用户名信息。
	Username *string `json:"username,omitempty"`

	// **参数说明**：转发kafka关联的密码信息。
	Password *string `json:"password,omitempty"`

	// **参数说明**：转发kafka关联的鉴权机制。 **取值范围**： - PAAS：非SASL鉴权。 - PLAIN：SASL/PLAIN模式。需要填写对应的用户名密码信息。
	Mechanism *string `json:"mechanism,omitempty"`
}

func (o DmsKafkaForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DmsKafkaForwarding struct{}"
	}

	return strings.Join([]string{"DmsKafkaForwarding", string(data)}, " ")
}
