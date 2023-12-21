package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPartitionRequest Request Object
type ShowPartitionRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 分区名称
	PartitionName string `json:"partition_name"`
}

func (o ShowPartitionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPartitionRequest struct{}"
	}

	return strings.Join([]string{"ShowPartitionRequest", string(data)}, " ")
}
