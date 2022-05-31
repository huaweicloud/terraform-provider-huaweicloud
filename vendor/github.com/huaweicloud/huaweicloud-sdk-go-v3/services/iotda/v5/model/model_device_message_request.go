package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeviceMessageRequest struct {

	// **参数说明**：消息id，由用户生成（推荐使用UUID），同一个设备下唯一， 如果用户填写的id在设备下不唯一， 则接口返回错误。 **取值范围**：长度不超过100，只允许字母、数字、下划线（_）、连接符（-）的组合。
	MessageId *string `json:"message_id,omitempty"`

	// **参数说明**：消息名称。 **取值范围**：长度不超过128，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	Name *string `json:"name,omitempty"`

	// **参数说明**：消息内容，支持string和json格式。
	Message *interface{} `json:"message"`

	// **参数说明**：消息内容编码格式。默认值none。 **取值范围**： - none  - base64：只能通过topic_full_name字段自定义的topic发送消息,否则会发送失败。
	Encoding *string `json:"encoding,omitempty"`

	// **参数说明**：有效负载格式，在消息内容编码格式为none时有效。默认值standard（平台封装的标准格式）。 **取值范围**： - standard  - raw：时直接将消息内容作为有效负载下发， 注意： 取值为raw时，只能通过topic_full_name字段自定义的topic发送消息，否则会发送失败。
	PayloadFormat *string `json:"payload_format,omitempty"`

	// **参数说明**：消息下行到设备的topic, 可选， 仅适用于MQTT协议接入的设备。 用户只能填写在租户产品界面配置的topic, 否则会校验不通过。 平台给消息topic添加的前缀为$oc/devices/{device_id}/user/， 用户可以在前缀的基础上增加自定义部分， 如增加messageDown，则平台拼接前缀后完整的topic为 $oc/devices/{device_id}/user/messageDown，其中device_id以实际设备的网关id替代。 如果用户指定该topic，消息会通过该topic下行到设备，如果用户不指定， 则消息通过系统默认的topic下行到设备,系统默认的topic格式为： $oc/devices/{device_id}/sys/messages/down。此字段与topic_full_name字段只可填写一个。
	Topic *string `json:"topic,omitempty"`

	// **参数说明**：消息下行到设备的完整topic名称, 可选。用户需要下发用户自定义的topic给设备时，可以使用该参数指定完整的topic名称，物联网平台不校验该topic是否在平台定义，直接透传给设备。设备需要提前订阅该topic。此字段与topic字段只可填写一个。
	TopicFullName *string `json:"topic_full_name,omitempty"`
}

func (o DeviceMessageRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeviceMessageRequest struct{}"
	}

	return strings.Join([]string{"DeviceMessageRequest", string(data)}, " ")
}
