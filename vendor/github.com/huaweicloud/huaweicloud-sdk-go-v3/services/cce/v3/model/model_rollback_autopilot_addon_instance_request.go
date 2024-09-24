package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RollbackAutopilotAddonInstanceRequest Request Object
type RollbackAutopilotAddonInstanceRequest struct {

	// 插件实例ID
	Id string `json:"id"`

	Body *AddonInstanceRollbackRequest `json:"body,omitempty"`
}

func (o RollbackAutopilotAddonInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RollbackAutopilotAddonInstanceRequest struct{}"
	}

	return strings.Join([]string{"RollbackAutopilotAddonInstanceRequest", string(data)}, " ")
}
