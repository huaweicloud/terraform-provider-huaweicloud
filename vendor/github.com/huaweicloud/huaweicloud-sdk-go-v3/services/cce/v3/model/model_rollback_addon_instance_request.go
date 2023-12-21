package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RollbackAddonInstanceRequest Request Object
type RollbackAddonInstanceRequest struct {

	// 插件实例ID
	Id string `json:"id"`

	Body *AddonInstanceRollbackRequest `json:"body,omitempty"`
}

func (o RollbackAddonInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RollbackAddonInstanceRequest struct{}"
	}

	return strings.Join([]string{"RollbackAddonInstanceRequest", string(data)}, " ")
}
