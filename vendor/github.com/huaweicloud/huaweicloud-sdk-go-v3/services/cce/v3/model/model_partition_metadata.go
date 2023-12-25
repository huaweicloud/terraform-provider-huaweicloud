package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PartitionMetadata 分区的元数据信息
type PartitionMetadata struct {

	// 分区名称
	Name *string `json:"name,omitempty"`

	// 创建时间
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`
}

func (o PartitionMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PartitionMetadata struct{}"
	}

	return strings.Join([]string{"PartitionMetadata", string(data)}, " ")
}
