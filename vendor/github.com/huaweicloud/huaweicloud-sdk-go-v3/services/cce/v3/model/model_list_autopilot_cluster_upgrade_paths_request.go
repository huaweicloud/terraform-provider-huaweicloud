package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotClusterUpgradePathsRequest Request Object
type ListAutopilotClusterUpgradePathsRequest struct {
}

func (o ListAutopilotClusterUpgradePathsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotClusterUpgradePathsRequest struct{}"
	}

	return strings.Join([]string{"ListAutopilotClusterUpgradePathsRequest", string(data)}, " ")
}
