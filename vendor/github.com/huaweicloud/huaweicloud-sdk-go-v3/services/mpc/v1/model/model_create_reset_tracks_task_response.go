package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateResetTracksTaskResponse struct {

	// 任务ID。 如果返回值为200 OK，为接受任务后产生的任务ID。
	TaskId         *string `json:"task_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateResetTracksTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateResetTracksTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateResetTracksTaskResponse", string(data)}, " ")
}
