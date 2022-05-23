package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateMessageRequest struct {

	// **参数说明**：下发消息的设备ID，用于唯一标识一个设备，在注册设备时由物联网平台分配获得。 **取值范围**：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。
	DeviceId string `json:"device_id"`

	// **参数说明**：实例ID。物理多租下各实例的唯一标识，一般华为云租户无需携带该参数，仅在物理多租场景下从管理面访问API时需要携带该参数。
	InstanceId *string `json:"Instance-Id,omitempty"`

	Body *DeviceMessageRequest `json:"body,omitempty"`
}

func (o CreateMessageRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMessageRequest struct{}"
	}

	return strings.Join([]string{"CreateMessageRequest", string(data)}, " ")
}
