package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutopilotCluster
type AutopilotCluster struct {

	// API类型，固定值“Cluster”或“cluster”，该值不可修改。
	Kind string `json:"kind"`

	// API版本，固定值“v3”，该值不可修改。
	ApiVersion string `json:"apiVersion"`

	Metadata *AutopilotClusterMetadata `json:"metadata"`

	Spec *AutopilotClusterSpec `json:"spec"`

	Status *AutopilotClusterStatus `json:"status,omitempty"`
}

func (o AutopilotCluster) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotCluster struct{}"
	}

	return strings.Join([]string{"AutopilotCluster", string(data)}, " ")
}
