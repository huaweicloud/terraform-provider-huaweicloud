package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTaskResponse Response Object
type ShowTaskResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	TaskInfo       *TaskInfo `json:"taskInfo,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ShowTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTaskResponse struct{}"
	}

	return strings.Join([]string{"ShowTaskResponse", string(data)}, " ")
}
