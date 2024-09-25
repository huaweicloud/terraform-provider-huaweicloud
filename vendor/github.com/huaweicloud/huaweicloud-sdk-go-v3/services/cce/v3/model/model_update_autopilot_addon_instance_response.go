package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAutopilotAddonInstanceResponse Response Object
type UpdateAutopilotAddonInstanceResponse struct {

	// API类型，固定值“Addon”，该值不可修改。
	Kind *string `json:"kind,omitempty"`

	// API版本，固定值“v3”，该值不可修改。
	ApiVersion *string `json:"apiVersion,omitempty"`

	Metadata *AddonMetadata `json:"metadata,omitempty"`

	Spec *InstanceSpec `json:"spec,omitempty"`

	Status         *AddonInstanceStatus `json:"status,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o UpdateAutopilotAddonInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAutopilotAddonInstanceResponse struct{}"
	}

	return strings.Join([]string{"UpdateAutopilotAddonInstanceResponse", string(data)}, " ")
}
