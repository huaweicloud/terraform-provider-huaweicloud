package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowDevicesInGroupResponse struct {

	// 设备列表。
	Devices *[]SimplifyDevice `json:"devices,omitempty"`

	Page           *Page `json:"page,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ShowDevicesInGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDevicesInGroupResponse struct{}"
	}

	return strings.Join([]string{"ShowDevicesInGroupResponse", string(data)}, " ")
}
