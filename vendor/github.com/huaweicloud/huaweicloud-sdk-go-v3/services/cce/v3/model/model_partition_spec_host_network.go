package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PartitionSpecHostNetwork 分区子网
type PartitionSpecHostNetwork struct {

	// 子网ID
	SubnetID *string `json:"subnetID,omitempty"`
}

func (o PartitionSpecHostNetwork) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PartitionSpecHostNetwork struct{}"
	}

	return strings.Join([]string{"PartitionSpecHostNetwork", string(data)}, " ")
}
