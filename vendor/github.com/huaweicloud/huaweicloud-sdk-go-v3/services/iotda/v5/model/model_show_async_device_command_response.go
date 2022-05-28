package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowAsyncDeviceCommandResponse struct {

	// 设备ID，用于唯一标识一个设备，在注册设备时由物联网平台分配获得。
	DeviceId *string `json:"device_id,omitempty"`

	// 设备命令ID，用于唯一标识一条命令，在下发设备命令时由物联网平台分配获得。
	CommandId *string `json:"command_id,omitempty"`

	// 设备命令所属的设备服务ID，在设备关联的产品模型中定义。
	ServiceId *string `json:"service_id,omitempty"`

	// 设备命令名称，在设备关联的产品模型中定义。
	CommandName *string `json:"command_name,omitempty"`

	// 设备执行的命令，Json格式，里面是一个个健值对，如果service_id不为空，每个健都是profile中命令的参数名（paraName）;如果service_id为空则由用户自定义命令格式。设备命令示例：{\"value\":\"1\"}，具体格式需要应用和设备约定。
	Paras *interface{} `json:"paras,omitempty"`

	// 物联网平台缓存命令的时长， 单位秒。
	ExpireTime *int32 `json:"expire_time,omitempty"`

	// 下发命令的状态。 ·PENDING表示未下发,在物联网平台缓存着 ·EXPIRED表示命令已经过期，即缓存的时间超过设定的expire_time ·SENT表示命令正在下发 ·DELIVERED表示命令已送达设备 ·SUCCESSFUL表示命令已经成功执行 ·FAILED表示命令执行失败 ·TIMEOUT表示命令下发之后，没有收到设备确认或者响应结果一定时间后超时
	Status *string `json:"status,omitempty"`

	// 设备命令执行的详细结果，由设备返回，Json格式。
	Result *interface{} `json:"result,omitempty"`

	// 命令的创建时间，\"yyyyMMdd'T'HHmmss'Z'\"格式的UTC字符串。
	CreatedTime *string `json:"created_time,omitempty"`

	// 物联网平台发送命令的时间，如果命令是立即下发， 则该时间与命令创建时间一致， 如果是缓存命令， 则是命令实际下发的时间。\"yyyyMMdd'T'HHmmss'Z'\"格式的UTC字符串。
	SentTime *string `json:"sent_time,omitempty"`

	// 物联网平台将命令送达到设备的时间，\"yyyyMMdd'T'HHmmss'Z'\"格式的UTC字符串
	DeliveredTime *string `json:"delivered_time,omitempty"`

	// 下发策略， immediately表示立即下发，delay表示缓存起来，等数据上报或者设备上线之后下发。
	SendStrategy *string `json:"send_strategy,omitempty"`

	// 设备响应命令的时间，\"yyyyMMdd'T'HHmmss'Z'\"格式的UTC字符串
	ResponseTime   *string `json:"response_time,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowAsyncDeviceCommandResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAsyncDeviceCommandResponse struct{}"
	}

	return strings.Join([]string{"ShowAsyncDeviceCommandResponse", string(data)}, " ")
}
