package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeviceBroadcastRequest struct {

	// **参数说明**：资源空间ID。此参数为非必选参数，存在多资源空间的用户需要使用该接口时，建议携带该参数指定广播消息所属的资源空间，否则广播消息将会向[[默认资源空间](https://support.huaweicloud.com/usermanual-iothub/iot_01_0006.html#section0)](tag:hws)[[默认资源空间](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_0006.html#section0)](tag:hws_hk)下订阅指定topic的设备发送。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AppId *string `json:"app_id,omitempty"`

	// **参数说明**：接收广播消息的完整Topic名称, 必选。用户需要发布广播消息给设备时，可以使用该参数指定完整的topic名称，物联网平台会向指定资源空间下订阅了该topic的所有在线设备发送消息。广播的topic无需控制台创建，但topic的前缀必须为$oc/broadcast/
	TopicFullName string `json:"topic_full_name"`

	// **参数说明**：广播消息的内容，用户需要将消息原文使用Base64编码，Base64编码后的长度不超过128KB。
	Message string `json:"message"`

	// **参数说明**：广播消息在平台缓存的老化时间，时间单位是分钟，默认值为0；ttl参数数值必须是5的倍数，即以5分钟为粒度；指定为0时表示不缓存消息，最大缓存时间1440分钟，即缓存一天；ttl>0时，一个topic订阅设备数限制为10，如果一个topic订阅设备数超过10，则接口返回错误。
	Ttl *int32 `json:"ttl,omitempty"`

	// **参数说明**：消息id，由用户生成（推荐使用UUID）。ttl> 0时，平台会缓存消息，需确保message_id是唯一的， 否则接口返回错误。 **取值范围**：长度不超过100，只允许字母、数字、下划线（_）、连接符（-）的组合。
	MessageId *string `json:"message_id,omitempty"`
}

func (o DeviceBroadcastRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeviceBroadcastRequest struct{}"
	}

	return strings.Join([]string{"DeviceBroadcastRequest", string(data)}, " ")
}
