package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateSubNetworkInterfaceRequest struct {

	// 辅助弹性网卡的唯一标识
	SubNetworkInterfaceId string `json:"sub_network_interface_id"`

	Body *UpdateSubNetworkInterfaceRequestBody `json:"body,omitempty"`
}

func (o UpdateSubNetworkInterfaceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSubNetworkInterfaceRequest struct{}"
	}

	return strings.Join([]string{"UpdateSubNetworkInterfaceRequest", string(data)}, " ")
}
