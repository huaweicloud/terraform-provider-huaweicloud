package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateInstanceRequest Request Object
type UpdateInstanceRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *UpdateInstanceReq `json:"body,omitempty"`
}

func (o UpdateInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceRequest struct{}"
	}

	return strings.Join([]string{"UpdateInstanceRequest", string(data)}, " ")
}
