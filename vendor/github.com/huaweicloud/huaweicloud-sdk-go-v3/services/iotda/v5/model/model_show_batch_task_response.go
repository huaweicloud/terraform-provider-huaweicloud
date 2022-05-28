package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowBatchTaskResponse struct {
	Batchtask *Task `json:"batchtask,omitempty"`

	// 子任务详情列表。
	TaskDetails *[]TaskDetail `json:"task_details,omitempty"`

	Page           *Page `json:"page,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ShowBatchTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBatchTaskResponse struct{}"
	}

	return strings.Join([]string{"ShowBatchTaskResponse", string(data)}, " ")
}
