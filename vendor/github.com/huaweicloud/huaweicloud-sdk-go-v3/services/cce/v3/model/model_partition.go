package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Partition 集群分区信息
type Partition struct {

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	Metadata *PartitionMetadata `json:"metadata,omitempty"`

	Spec *PartitionSpec `json:"spec,omitempty"`
}

func (o Partition) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Partition struct{}"
	}

	return strings.Join([]string{"Partition", string(data)}, " ")
}
