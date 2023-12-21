package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PartitionReqBody 集群分区信息
type PartitionReqBody struct {

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	Metadata *PartitionReqBodyMetadata `json:"metadata,omitempty"`

	Spec *PartitionSpec `json:"spec,omitempty"`
}

func (o PartitionReqBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PartitionReqBody struct{}"
	}

	return strings.Join([]string{"PartitionReqBody", string(data)}, " ")
}
