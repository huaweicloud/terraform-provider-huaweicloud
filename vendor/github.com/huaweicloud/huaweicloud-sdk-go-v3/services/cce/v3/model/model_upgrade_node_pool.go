package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeNodePool
type UpgradeNodePool struct {

	// API类型，固定值“NodePool”。
	Kind *string `json:"kind,omitempty"`

	// API版本，固定值“v3”。
	ApiVersion *string `json:"apiVersion,omitempty"`

	Metadata *NodePoolMetadata `json:"metadata,omitempty"`

	Spec *NodePoolUpgradeSpec `json:"spec"`
}

func (o UpgradeNodePool) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeNodePool struct{}"
	}

	return strings.Join([]string{"UpgradeNodePool", string(data)}, " ")
}
