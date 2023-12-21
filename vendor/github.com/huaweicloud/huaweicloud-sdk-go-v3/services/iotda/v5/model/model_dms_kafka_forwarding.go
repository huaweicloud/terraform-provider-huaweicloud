package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DmsKafkaForwarding 转发kafka消息结构
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

	// **参数说明**：转发kafka关联的SASL认证机制。 **取值范围**： - PAAS：明文传输，此模式下为非数据加密传输模式，数据传输不安全，建议您使用更安全的数据加密模式。 - PLAIN：SASL/PLAIN模式。需要填写对应的用户名密码信息。一种简单的用户名密码校验机制，在SASL_PLAINTEXT场景下，不建议使用。 - SCRAM-SHA-512：SASL/SCRAM-SHA-512模式。需要填写对应的用户名密码信息。采用哈希算法对用户名与密码生成凭证，进行身份校验的安全认证机制，比PLAIN机制安全性更高。
	Mechanism *string `json:"mechanism,omitempty"`

	// **参数说明**：kafka传输安全协议，此字段不填默认为SASL_SSL。当mechanism为PAAS或不填时，该字段不生效。 **取值范围**： - SASL_SSL：采用SSL证书进行加密传输，支持帐号密码认证，安全性更高。 - SASL_PLAINTEXT：明文传输，支持帐号密码认证，性能更好，建议mechanism使用SCRAM-SHA-512机制。
	SecurityProtocol *string `json:"security_protocol,omitempty"`
}

func (o DmsKafkaForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DmsKafkaForwarding struct{}"
	}

	return strings.Join([]string{"DmsKafkaForwarding", string(data)}, " ")
}
