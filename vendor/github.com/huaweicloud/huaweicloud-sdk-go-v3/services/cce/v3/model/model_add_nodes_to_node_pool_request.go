package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddNodesToNodePoolRequest Request Object
type AddNodesToNodePoolRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 节点池ID
	NodepoolId string `json:"nodepool_id"`

	Body *AddNodesToNodePoolList `json:"body,omitempty"`
}

func (o AddNodesToNodePoolRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddNodesToNodePoolRequest struct{}"
	}

	return strings.Join([]string{"AddNodesToNodePoolRequest", string(data)}, " ")
}
