package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateAsyncCommandResponse struct {

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

	// 设备命令状态,如果命令被缓存，返回PENDING, 如果命令下发给设备，返回SENT。
	Status *string `json:"status,omitempty"`

	// 命令的创建时间，\"yyyyMMdd'T'HHmmss'Z'\"格式的UTC字符串。
	CreatedTime *string `json:"created_time,omitempty"`

	// 下发策略， immediately表示立即下发，delay表示缓存起来，等数据上报或者设备上线之后下发。
	SendStrategy   *string `json:"send_strategy,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateAsyncCommandResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAsyncCommandResponse struct{}"
	}

	return strings.Join([]string{"CreateAsyncCommandResponse", string(data)}, " ")
}
