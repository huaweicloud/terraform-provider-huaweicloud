package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeNodePoolRequest Request Object
type UpgradeNodePoolRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 节点池ID
	NodepoolId string `json:"nodepool_id"`

	Body *UpgradeNodePool `json:"body,omitempty"`
}

func (o UpgradeNodePoolRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeNodePoolRequest struct{}"
	}

	return strings.Join([]string{"UpgradeNodePoolRequest", string(data)}, " ")
}
