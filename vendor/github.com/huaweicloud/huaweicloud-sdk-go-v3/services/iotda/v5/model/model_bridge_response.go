package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BridgeResponse struct {

	// 网桥ID
	BridgeId *string `json:"bridge_id,omitempty"`

	// 网桥名称。
	BridgeName *string `json:"bridge_name,omitempty"`

	// 网桥状态。 - ONLINE：网桥在线。 - OFFLINE：网桥离线。
	Status *string `json:"status,omitempty"`
}

func (o BridgeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BridgeResponse struct{}"
	}

	return strings.Join([]string{"BridgeResponse", string(data)}, " ")
}
