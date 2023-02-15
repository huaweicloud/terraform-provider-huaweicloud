package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 端侧设备信息
type DeviceSide struct {

	// **参数说明**：端侧执行下发的目标设备ID列表。设备ID，用于唯一标识一个设备，在注册设备时由物联网平台分配获得。
	DeviceIds *[]string `json:"device_ids,omitempty"`
}

func (o DeviceSide) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeviceSide struct{}"
	}

	return strings.Join([]string{"DeviceSide", string(data)}, " ")
}
