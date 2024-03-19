package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTaskRelatedTestCaseResponse Response Object
type UpdateTaskRelatedTestCaseResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	TaskInfo       *TaskInfo `json:"taskInfo,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o UpdateTaskRelatedTestCaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskRelatedTestCaseResponse struct{}"
	}

	return strings.Join([]string{"UpdateTaskRelatedTestCaseResponse", string(data)}, " ")
}
