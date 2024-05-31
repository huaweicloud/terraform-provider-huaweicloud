package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDeviceGroupsByDeviceResponse Response Object
type ListDeviceGroupsByDeviceResponse struct {

	// 设备组信息列表。
	DeviceGroups   *[]ListDeviceGroupSummary `json:"device_groups,omitempty"`
	HttpStatusCode int                       `json:"-"`
}

func (o ListDeviceGroupsByDeviceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDeviceGroupsByDeviceResponse struct{}"
	}

	return strings.Join([]string{"ListDeviceGroupsByDeviceResponse", string(data)}, " ")
}
