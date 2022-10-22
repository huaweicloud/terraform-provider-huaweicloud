package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type StartTaskGroupRequest struct {

	// 任务组id
	GroupId string `json:"group_id"`

	Body *StartTaskGroupReq `json:"body,omitempty"`
}

func (o StartTaskGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartTaskGroupRequest struct{}"
	}

	return strings.Join([]string{"StartTaskGroupRequest", string(data)}, " ")
}
