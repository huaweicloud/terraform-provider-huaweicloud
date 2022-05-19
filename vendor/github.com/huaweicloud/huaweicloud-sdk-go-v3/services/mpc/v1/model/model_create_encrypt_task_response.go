package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateEncryptTaskResponse struct {

	// 加密任务Id
	TaskId         *string `json:"task_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateEncryptTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateEncryptTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateEncryptTaskResponse", string(data)}, " ")
}
