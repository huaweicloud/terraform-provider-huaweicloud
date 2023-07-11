package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClusterRequest Request Object
type ShowClusterRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o ShowClusterRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterRequest struct{}"
	}

	return strings.Join([]string{"ShowClusterRequest", string(data)}, " ")
}
