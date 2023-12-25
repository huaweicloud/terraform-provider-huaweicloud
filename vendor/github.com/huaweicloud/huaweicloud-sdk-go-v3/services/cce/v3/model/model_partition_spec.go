package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PartitionSpec 分区的配置信息
type PartitionSpec struct {
	HostNetwork *PartitionSpecHostNetwork `json:"hostNetwork,omitempty"`

	// 分区容器子网
	ContainerNetwork *[]PartitionSpecContainerNetwork `json:"containerNetwork,omitempty"`

	// 群组
	PublicBorderGroup *string `json:"publicBorderGroup,omitempty"`

	// 类别
	Category *string `json:"category,omitempty"`
}

func (o PartitionSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PartitionSpec struct{}"
	}

	return strings.Join([]string{"PartitionSpec", string(data)}, " ")
}
