package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RetryUpgradeTaskRequest Request Object
type RetryUpgradeTaskRequest struct {

	// 待重试的集群ID。
	ClusterId string `json:"cluster_id"`

	// 待重试的任务ID。
	ActionId string `json:"action_id"`

	// 当该参数不为空时，终止该任务的影响。当前仅支持abort。
	RetryMode *string `json:"retry_mode,omitempty"`
}

func (o RetryUpgradeTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RetryUpgradeTaskRequest struct{}"
	}

	return strings.Join([]string{"RetryUpgradeTaskRequest", string(data)}, " ")
}
