package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddTunnelDto 创建隧道请求结构体
type AddTunnelDto struct {

	// **参数说明**：设备ID
	DeviceId string `json:"device_id"`
}

func (o AddTunnelDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddTunnelDto struct{}"
	}

	return strings.Join([]string{"AddTunnelDto", string(data)}, " ")
}
