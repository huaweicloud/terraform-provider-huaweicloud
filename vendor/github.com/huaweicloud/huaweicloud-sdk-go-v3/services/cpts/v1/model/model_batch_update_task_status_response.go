package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchUpdateTaskStatusResponse Response Object
type BatchUpdateTaskStatusResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 扩展字段
	Extend *interface{} `json:"extend,omitempty"`

	// 批量启停任务响应结果
	Result         *[]BatchUpdateTaskStatusResponseBodyResult `json:"result,omitempty"`
	HttpStatusCode int                                        `json:"-"`
}

func (o BatchUpdateTaskStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchUpdateTaskStatusResponse struct{}"
	}

	return strings.Join([]string{"BatchUpdateTaskStatusResponse", string(data)}, " ")
}
