package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateAutopilotClusterRequest Request Object
type CreateAutopilotClusterRequest struct {
	Body *AutopilotCluster `json:"body,omitempty"`
}

func (o CreateAutopilotClusterRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAutopilotClusterRequest struct{}"
	}

	return strings.Join([]string{"CreateAutopilotClusterRequest", string(data)}, " ")
}
