package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPreCheckResponse Response Object
type ShowPreCheckResponse struct {

	// api版本，默认为v3
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 资源类型，默认为PreCheckTask
	Kind *string `json:"kind,omitempty"`

	Metadata *PrecheckTaskMetadata `json:"metadata,omitempty"`

	Spec *PrecheckSpec `json:"spec,omitempty"`

	Status         *PrecheckStatus `json:"status,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ShowPreCheckResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPreCheckResponse struct{}"
	}

	return strings.Join([]string{"ShowPreCheckResponse", string(data)}, " ")
}
