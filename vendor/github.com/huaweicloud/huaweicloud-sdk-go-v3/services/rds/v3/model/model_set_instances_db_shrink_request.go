package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetInstancesDbShrinkRequest Request Object
type SetInstancesDbShrinkRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *UpdateDbShrinkRequestBody `json:"body,omitempty"`
}

func (o SetInstancesDbShrinkRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetInstancesDbShrinkRequest struct{}"
	}

	return strings.Join([]string{"SetInstancesDbShrinkRequest", string(data)}, " ")
}
