package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutoCreatePolicyRequest Request Object
type ShowAutoCreatePolicyRequest struct {

	// 指定需查询自动创建快照策略的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ShowAutoCreatePolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutoCreatePolicyRequest struct{}"
	}

	return strings.Join([]string{"ShowAutoCreatePolicyRequest", string(data)}, " ")
}
