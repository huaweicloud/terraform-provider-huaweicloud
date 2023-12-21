package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PartitionSpecContainerNetwork struct {

	// 子网ID
	SubnetID *string `json:"subnetID,omitempty"`
}

func (o PartitionSpecContainerNetwork) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PartitionSpecContainerNetwork struct{}"
	}

	return strings.Join([]string{"PartitionSpecContainerNetwork", string(data)}, " ")
}
