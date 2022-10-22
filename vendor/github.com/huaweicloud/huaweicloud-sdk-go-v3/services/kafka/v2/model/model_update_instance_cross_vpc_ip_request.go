package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateInstanceCrossVpcIpRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *UpdateInstanceCrossVpcIpReq `json:"body,omitempty"`
}

func (o UpdateInstanceCrossVpcIpRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceCrossVpcIpRequest struct{}"
	}

	return strings.Join([]string{"UpdateInstanceCrossVpcIpRequest", string(data)}, " ")
}
