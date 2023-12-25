package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClusterEndpointsRequest Request Object
type ShowClusterEndpointsRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`
}

func (o ShowClusterEndpointsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterEndpointsRequest struct{}"
	}

	return strings.Join([]string{"ShowClusterEndpointsRequest", string(data)}, " ")
}
