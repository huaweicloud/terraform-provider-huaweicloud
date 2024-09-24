package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotClusterUpgradeFeatureGatesRequest Request Object
type ListAutopilotClusterUpgradeFeatureGatesRequest struct {
}

func (o ListAutopilotClusterUpgradeFeatureGatesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotClusterUpgradeFeatureGatesRequest struct{}"
	}

	return strings.Join([]string{"ListAutopilotClusterUpgradeFeatureGatesRequest", string(data)}, " ")
}
