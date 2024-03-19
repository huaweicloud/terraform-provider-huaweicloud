package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListUpgradeWorkFlowsRequest Request Object
type ListUpgradeWorkFlowsRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`
}

func (o ListUpgradeWorkFlowsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUpgradeWorkFlowsRequest struct{}"
	}

	return strings.Join([]string{"ListUpgradeWorkFlowsRequest", string(data)}, " ")
}
