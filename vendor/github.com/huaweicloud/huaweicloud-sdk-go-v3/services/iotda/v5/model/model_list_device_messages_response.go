package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListDeviceMessagesResponse struct {

	// 设备ID，用于唯一标识一个设备，在注册设备时由物联网平台分配获得。
	DeviceId *string `json:"device_id,omitempty"`

	// 设备消息列表。
	Messages       *[]DeviceMessage `json:"messages,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListDeviceMessagesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDeviceMessagesResponse struct{}"
	}

	return strings.Join([]string{"ListDeviceMessagesResponse", string(data)}, " ")
}
