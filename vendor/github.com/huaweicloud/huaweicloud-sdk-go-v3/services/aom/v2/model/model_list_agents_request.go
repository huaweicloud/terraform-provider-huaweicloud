package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAgentsRequest Request Object
type ListAgentsRequest struct {

	// 集群id。
	ClusterId string `json:"cluster_id"`

	// 命名空间。
	Namespace string `json:"namespace"`
}

func (o ListAgentsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAgentsRequest struct{}"
	}

	return strings.Join([]string{"ListAgentsRequest", string(data)}, " ")
}
