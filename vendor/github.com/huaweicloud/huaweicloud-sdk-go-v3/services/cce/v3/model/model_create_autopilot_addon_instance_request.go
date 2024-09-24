package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateAutopilotAddonInstanceRequest Request Object
type CreateAutopilotAddonInstanceRequest struct {
	Body *InstanceRequest `json:"body,omitempty"`
}

func (o CreateAutopilotAddonInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAutopilotAddonInstanceRequest struct{}"
	}

	return strings.Join([]string{"CreateAutopilotAddonInstanceRequest", string(data)}, " ")
}
