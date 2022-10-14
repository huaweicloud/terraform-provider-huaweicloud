package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTaskGroupResponse struct {

	// 创建的迁移任务组id
	GroupId        *string `json:"group_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateTaskGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTaskGroupResponse struct{}"
	}

	return strings.Join([]string{"CreateTaskGroupResponse", string(data)}, " ")
}
