package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutopilotPreCheckResponse Response Object
type ShowAutopilotPreCheckResponse struct {

	// api版本，默认为v3
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 资源类型，默认为PreCheckTask
	Kind *string `json:"kind,omitempty"`

	Metadata *PrecheckTaskMetadata `json:"metadata,omitempty"`

	Spec *PrecheckSpec `json:"spec,omitempty"`

	Status         *PrecheckStatus `json:"status,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ShowAutopilotPreCheckResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutopilotPreCheckResponse struct{}"
	}

	return strings.Join([]string{"ShowAutopilotPreCheckResponse", string(data)}, " ")
}
