package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTaskGroupRequest Request Object
type UpdateTaskGroupRequest struct {

	// 任务组id
	GroupId string `json:"group_id"`

	Body *UpdateBandwidthPolicyReq `json:"body,omitempty"`
}

func (o UpdateTaskGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskGroupRequest struct{}"
	}

	return strings.Join([]string{"UpdateTaskGroupRequest", string(data)}, " ")
}
