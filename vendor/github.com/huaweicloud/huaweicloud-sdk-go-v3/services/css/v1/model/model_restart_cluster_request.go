package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type RestartClusterRequest struct {

	// 指定重启集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o RestartClusterRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestartClusterRequest struct{}"
	}

	return strings.Join([]string{"RestartClusterRequest", string(data)}, " ")
}
