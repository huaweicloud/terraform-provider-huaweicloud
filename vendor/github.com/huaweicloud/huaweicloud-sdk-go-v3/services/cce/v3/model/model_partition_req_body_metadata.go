package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PartitionReqBodyMetadata 分区的元数据信息
type PartitionReqBodyMetadata struct {

	// 分区名称
	Name *string `json:"name,omitempty"`
}

func (o PartitionReqBodyMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PartitionReqBodyMetadata struct{}"
	}

	return strings.Join([]string{"PartitionReqBodyMetadata", string(data)}, " ")
}
