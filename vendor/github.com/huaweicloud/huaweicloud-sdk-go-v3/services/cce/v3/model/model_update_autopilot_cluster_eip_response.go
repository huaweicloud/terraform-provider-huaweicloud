package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAutopilotClusterEipResponse Response Object
type UpdateAutopilotClusterEipResponse struct {
	Metadata *Metadata `json:"metadata,omitempty"`

	Spec *MasterEipResponseSpec `json:"spec,omitempty"`

	Status         *MasterEipResponseStatus `json:"status,omitempty"`
	HttpStatusCode int                      `json:"-"`
}

func (o UpdateAutopilotClusterEipResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAutopilotClusterEipResponse struct{}"
	}

	return strings.Join([]string{"UpdateAutopilotClusterEipResponse", string(data)}, " ")
}
