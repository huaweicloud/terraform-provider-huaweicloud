package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartTargetClusterConnectivityTestRequest Request Object
type StartTargetClusterConnectivityTestRequest struct {

	// 指定待测试的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *StartTargetClusterConnectivityTestReq `json:"body,omitempty"`
}

func (o StartTargetClusterConnectivityTestRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartTargetClusterConnectivityTestRequest struct{}"
	}

	return strings.Join([]string{"StartTargetClusterConnectivityTestRequest", string(data)}, " ")
}
