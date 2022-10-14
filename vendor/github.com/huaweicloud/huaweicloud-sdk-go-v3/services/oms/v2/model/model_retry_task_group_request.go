package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type RetryTaskGroupRequest struct {

	// 任务组id
	GroupId string `json:"group_id"`

	Body *RetryTaskGroupReq `json:"body,omitempty"`
}

func (o RetryTaskGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RetryTaskGroupRequest struct{}"
	}

	return strings.Join([]string{"RetryTaskGroupRequest", string(data)}, " ")
}
