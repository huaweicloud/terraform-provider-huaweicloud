package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePreCheckResponse Response Object
type CreatePreCheckResponse struct {

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	Metadata *PrecheckCluserResponseMetadata `json:"metadata,omitempty"`

	Spec *PrecheckSpec `json:"spec,omitempty"`

	Status         *PrecheckStatus `json:"status,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o CreatePreCheckResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePreCheckResponse struct{}"
	}

	return strings.Join([]string{"CreatePreCheckResponse", string(data)}, " ")
}
