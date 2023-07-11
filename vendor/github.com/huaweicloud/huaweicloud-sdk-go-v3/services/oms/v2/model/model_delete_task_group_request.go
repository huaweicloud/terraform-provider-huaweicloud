package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteTaskGroupRequest Request Object
type DeleteTaskGroupRequest struct {

	// 任务组id
	GroupId string `json:"group_id"`
}

func (o DeleteTaskGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTaskGroupRequest struct{}"
	}

	return strings.Join([]string{"DeleteTaskGroupRequest", string(data)}, " ")
}
