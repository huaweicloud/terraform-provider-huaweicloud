package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeviceMessage struct {

	// 设备消息ID，用于唯一标识一条消息，在下发设备消息时由物联网平台分配获得。
	MessageId *string `json:"message_id,omitempty"`

	// 消息名称,在下发消息时由用户指定。
	Name *string `json:"name,omitempty"`

	// 消息内容。
	Message *interface{} `json:"message,omitempty"`

	// 消息内容编码格式，取值范围none|base64,默认值none, base64格式仅支持透传。
	Encoding *string `json:"encoding,omitempty"`

	// 有效负载格式，在消息内容编码格式为none时有效，取值范围standard|raw，默认值standard（平台封装的标准格式），取值为raw时直接将消息内容作为有效负载下发。
	PayloadFormat *string `json:"payload_format,omitempty"`

	// 消息topic
	Topic *string `json:"topic,omitempty"`

	// 消息状态，包含PENDING，DELIVERED，FAILED和TIMEOUT，PENDING指设备不在线，消息被缓存起来，等设备上线之后下发； DELIVERED指消息发送成功；FAILED消息发送失败；TIMEOUT指消息在平台默认时间内（1天）还没有下发送给设备，则平台会将消息设置为超时，状态为TIMEOUT。
	Status *string `json:"status,omitempty"`

	// 消息的创建时间，\"yyyyMMdd'T'HHmmss'Z'\"格式的UTC字符串。
	CreatedTime *string `json:"created_time,omitempty"`

	// 消息结束时间, \"yyyyMMdd'T'HHmmss'Z'\"格式的UTC字符串，包含消息转换到DELIVERED和TIMEOUT两个状态的时间。
	FinishedTime *string `json:"finished_time,omitempty"`
}

func (o DeviceMessage) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeviceMessage struct{}"
	}

	return strings.Join([]string{"DeviceMessage", string(data)}, " ")
}
