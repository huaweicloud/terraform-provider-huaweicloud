package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeSpecUpdateNodeNicSpecUpdate 更新节点的网卡信息
type NodeSpecUpdateNodeNicSpecUpdate struct {
	PrimaryNic *NodeSpecUpdateNodeNicSpecUpdatePrimaryNic `json:"primaryNic,omitempty"`
}

func (o NodeSpecUpdateNodeNicSpecUpdate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeSpecUpdateNodeNicSpecUpdate struct{}"
	}

	return strings.Join([]string{"NodeSpecUpdateNodeNicSpecUpdate", string(data)}, " ")
}
