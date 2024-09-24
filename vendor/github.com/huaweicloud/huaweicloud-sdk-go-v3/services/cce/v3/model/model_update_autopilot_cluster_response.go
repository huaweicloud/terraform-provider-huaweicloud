package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAutopilotClusterResponse Response Object
type UpdateAutopilotClusterResponse struct {

	// API类型，固定值“Cluster”或“cluster”，该值不可修改。
	Kind *string `json:"kind,omitempty"`

	// API版本，固定值“v3”，该值不可修改。
	ApiVersion *string `json:"apiVersion,omitempty"`

	Metadata *AutopilotClusterMetadata `json:"metadata,omitempty"`

	Spec *AutopilotClusterSpec `json:"spec,omitempty"`

	Status         *AutopilotClusterStatus `json:"status,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o UpdateAutopilotClusterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAutopilotClusterResponse struct{}"
	}

	return strings.Join([]string{"UpdateAutopilotClusterResponse", string(data)}, " ")
}
