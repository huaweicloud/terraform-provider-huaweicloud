package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePartitionResponse Response Object
type CreatePartitionResponse struct {

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	Metadata *PartitionMetadata `json:"metadata,omitempty"`

	Spec           *PartitionSpec `json:"spec,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o CreatePartitionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePartitionResponse struct{}"
	}

	return strings.Join([]string{"CreatePartitionResponse", string(data)}, " ")
}
