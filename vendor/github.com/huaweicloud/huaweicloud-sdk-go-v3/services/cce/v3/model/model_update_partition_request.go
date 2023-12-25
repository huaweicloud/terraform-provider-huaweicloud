package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePartitionRequest Request Object
type UpdatePartitionRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 分区名称
	PartitionName string `json:"partition_name"`

	Body *PartitionReqBody `json:"body,omitempty"`
}

func (o UpdatePartitionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePartitionRequest struct{}"
	}

	return strings.Join([]string{"UpdatePartitionRequest", string(data)}, " ")
}
