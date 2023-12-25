package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPartitionResponse Response Object
type ShowPartitionResponse struct {

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	Metadata *PartitionMetadata `json:"metadata,omitempty"`

	Spec           *PartitionSpec `json:"spec,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ShowPartitionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPartitionResponse struct{}"
	}

	return strings.Join([]string{"ShowPartitionResponse", string(data)}, " ")
}
