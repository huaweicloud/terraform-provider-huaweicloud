package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodePoolUpdate
type NodePoolUpdate struct {
	Metadata *NodePoolMetadataUpdate `json:"metadata"`

	Spec *NodePoolSpecUpdate `json:"spec"`
}

func (o NodePoolUpdate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodePoolUpdate struct{}"
	}

	return strings.Join([]string{"NodePoolUpdate", string(data)}, " ")
}
