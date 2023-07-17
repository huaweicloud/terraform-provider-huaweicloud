package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTaskGroupRequest Request Object
type ShowTaskGroupRequest struct {

	// 任务组id
	GroupId string `json:"group_id"`
}

func (o ShowTaskGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTaskGroupRequest struct{}"
	}

	return strings.Join([]string{"ShowTaskGroupRequest", string(data)}, " ")
}
