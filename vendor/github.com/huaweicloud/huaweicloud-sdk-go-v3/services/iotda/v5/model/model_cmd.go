package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 给设备下发命令的命令内容，在动作的type为“DEVICE_CMD”时有效，且为必选
type Cmd struct {

	// **参数说明**：设备命令名称，在设备关联的产品模型中定义。
	CommandName string `json:"command_name"`

	// **参数说明**：设备命令参数，Json格式。 使用LWM2M协议设备命令示例：{\"value\":\"1\"}，里面是一个个健值对，每个健都是产品模型中命令的参数名（paraName）。 使用MQTT协议设备命令示例：{\"header\": {\"mode\": \"ACK\",\"from\": \"/users/testUser\",\"method\": \"SET_TEMPERATURE_READ_PERIOD\",\"to\":\"/devices/{device_id}/services/{service_id}\"},\"body\": {\"value\" : \"1\"}}。 - mode：必选，设备收到命令后是否需要回复确认消息，默认为ACK模式。ACK表示需要回复确认消息，NOACK表示不需要回复确认消息，其它值无效。 - from：可选，命令发送方的地址。App发起请求时格式为/users/{userId} ，应用服务器发起请求时格式为/{serviceName}，物联网平台发起请求时格式为/cloud/{serviceName}。 - to：可选，命令接收方的地址，格式为/devices/{device_id}/services/{service_id}。 - method：可选，产品模型中定义的命令名称。 - body：可选，命令的消息体，里面是一个个键值对，每个键都是产品模型中命令的参数名（paraName）。具体格式需要应用和设备约定。
	CommandBody *interface{} `json:"command_body"`

	// **参数说明**：设备命令所属的设备服务ID，在设备关联的产品模型中定义。
	ServiceId string `json:"service_id"`

	// **参数说明**：设备命令的缓存时间，单位为秒，表示物联网平台在把命令下发给设备前缓存命令的有效时间，超过这个时间后命令将不再下发，默认值为172800s（48小时）。如果buffer_timeout设置为0，则无论物联网平台上设置的命令下发模式是什么，该命令都会立即下发给设备。
	BufferTimeout *int32 `json:"buffer_timeout,omitempty"`

	// **参数说明**：命令响应的有效时间，单位为秒，表示设备接收到命令后，在response_timeout时间内响应有效，超过这个时间未收到命令的响应，则认为命令响应超时，默认值为1800s。
	ResponseTimeout *int32 `json:"response_timeout,omitempty"`

	// **参数说明**：设备命令的下发模式，仅当buffer_timeout的值大于0时有效。  - ACTIVE：主动模式，物联网平台主动将命令下发给设备。 - PASSIVE：被动模式，物联网平台创建设备命令后，会直接缓存命令。等到设备再次上线或者上报上一条命令的执行结果后才下发命令。
	Mode *string `json:"mode,omitempty"`
}

func (o Cmd) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Cmd struct{}"
	}

	return strings.Join([]string{"Cmd", string(data)}, " ")
}
