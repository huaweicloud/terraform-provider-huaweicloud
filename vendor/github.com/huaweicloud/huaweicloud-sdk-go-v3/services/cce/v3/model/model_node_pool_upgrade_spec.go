package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type NodePoolUpgradeSpec struct {

	// 节点池id。
	NodePoolID string `json:"nodePoolID"`

	NodeIDs *[]string `json:"nodeIDs,omitempty"`

	// Pod无法驱逐时，是否强制重置。
	Force *bool `json:"force,omitempty"`

	NodeTemplate *NodeTemplate `json:"nodeTemplate,omitempty"`

	MaxUnavailable *int32 `json:"maxUnavailable,omitempty"`

	RetryTimes *int32 `json:"retryTimes,omitempty"`

	SkippedNodes *[]string `json:"skippedNodes,omitempty"`
}

func (o NodePoolUpgradeSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodePoolUpgradeSpec struct{}"
	}

	return strings.Join([]string{"NodePoolUpgradeSpec", string(data)}, " ")
}
