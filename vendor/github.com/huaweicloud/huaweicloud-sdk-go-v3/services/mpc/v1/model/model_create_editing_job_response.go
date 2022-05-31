package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateEditingJobResponse struct {

	// 接受任务后产生的任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateEditingJobResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateEditingJobResponse struct{}"
	}

	return strings.Join([]string{"CreateEditingJobResponse", string(data)}, " ")
}
