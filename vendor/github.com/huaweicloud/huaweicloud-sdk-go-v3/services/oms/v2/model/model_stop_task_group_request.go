package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopTaskGroupRequest Request Object
type StopTaskGroupRequest struct {

	// 任务组id
	GroupId string `json:"group_id"`
}

func (o StopTaskGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopTaskGroupRequest struct{}"
	}

	return strings.Join([]string{"StopTaskGroupRequest", string(data)}, " ")
}
