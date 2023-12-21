package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteNodePoolRequest Request Object
type DeleteNodePoolRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 节点池ID
	NodepoolId string `json:"nodepool_id"`

	// 集群状态兼容Error参数，用于API平滑切换。 兼容场景下，errorStatus为空则屏蔽Error状态为Deleting状态。
	ErrorStatus *string `json:"errorStatus,omitempty"`
}

func (o DeleteNodePoolRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteNodePoolRequest struct{}"
	}

	return strings.Join([]string{"DeleteNodePoolRequest", string(data)}, " ")
}
