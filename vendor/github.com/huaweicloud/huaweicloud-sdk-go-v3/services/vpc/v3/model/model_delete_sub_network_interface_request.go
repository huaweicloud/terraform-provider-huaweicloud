package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteSubNetworkInterfaceRequest struct {

	// 弹性辅助网卡唯一标识
	SubNetworkInterfaceId string `json:"sub_network_interface_id"`
}

func (o DeleteSubNetworkInterfaceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSubNetworkInterfaceRequest struct{}"
	}

	return strings.Join([]string{"DeleteSubNetworkInterfaceRequest", string(data)}, " ")
}
