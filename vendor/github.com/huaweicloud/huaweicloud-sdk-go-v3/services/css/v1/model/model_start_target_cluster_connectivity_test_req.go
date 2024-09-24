package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartTargetClusterConnectivityTestReq struct {

	// 目标集群ID
	TargetClusterId string `json:"target_cluster_id"`
}

func (o StartTargetClusterConnectivityTestReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartTargetClusterConnectivityTestReq struct{}"
	}

	return strings.Join([]string{"StartTargetClusterConnectivityTestReq", string(data)}, " ")
}
