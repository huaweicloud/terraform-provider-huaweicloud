package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAgentHealthStatusRequest Request Object
type UpdateAgentHealthStatusRequest struct {

	// 探针id
	AgentId int32 `json:"agent_id"`

	Body *UpdateAgentHealthStatusRequestBody `json:"body,omitempty"`
}

func (o UpdateAgentHealthStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAgentHealthStatusRequest struct{}"
	}

	return strings.Join([]string{"UpdateAgentHealthStatusRequest", string(data)}, " ")
}
