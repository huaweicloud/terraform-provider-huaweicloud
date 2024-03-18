package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetInstancesNewDbShrinkRequest Request Object
type SetInstancesNewDbShrinkRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *UpdateDbShrinkRequestBody `json:"body,omitempty"`
}

func (o SetInstancesNewDbShrinkRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetInstancesNewDbShrinkRequest struct{}"
	}

	return strings.Join([]string{"SetInstancesNewDbShrinkRequest", string(data)}, " ")
}
