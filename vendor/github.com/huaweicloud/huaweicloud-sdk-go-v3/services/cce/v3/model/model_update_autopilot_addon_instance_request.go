package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAutopilotAddonInstanceRequest Request Object
type UpdateAutopilotAddonInstanceRequest struct {

	// 插件实例id
	Id string `json:"id"`

	Body *InstanceRequest `json:"body,omitempty"`
}

func (o UpdateAutopilotAddonInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAutopilotAddonInstanceRequest struct{}"
	}

	return strings.Join([]string{"UpdateAutopilotAddonInstanceRequest", string(data)}, " ")
}
