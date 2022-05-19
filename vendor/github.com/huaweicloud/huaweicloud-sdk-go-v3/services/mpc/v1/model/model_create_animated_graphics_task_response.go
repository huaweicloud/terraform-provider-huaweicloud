package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateAnimatedGraphicsTaskResponse struct {

	// 任务ID
	TaskId         *string `json:"task_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateAnimatedGraphicsTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAnimatedGraphicsTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateAnimatedGraphicsTaskResponse", string(data)}, " ")
}
