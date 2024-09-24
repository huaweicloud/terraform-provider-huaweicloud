package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutopilotClusterEndpointsResponse Response Object
type ShowAutopilotClusterEndpointsResponse struct {
	Metadata *Metadata `json:"metadata,omitempty"`

	Spec *OpenApiSpec `json:"spec,omitempty"`

	Status         *MasterEipResponseStatus `json:"status,omitempty"`
	HttpStatusCode int                      `json:"-"`
}

func (o ShowAutopilotClusterEndpointsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutopilotClusterEndpointsResponse struct{}"
	}

	return strings.Join([]string{"ShowAutopilotClusterEndpointsResponse", string(data)}, " ")
}
