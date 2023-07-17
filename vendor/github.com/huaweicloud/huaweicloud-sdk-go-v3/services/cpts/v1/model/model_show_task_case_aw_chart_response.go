package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTaskCaseAwChartResponse Response Object
type ShowTaskCaseAwChartResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 扩展字段
	Extend *interface{} `json:"extend,omitempty"`

	Result         *TaskCaseAwChartResult `json:"result,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

func (o ShowTaskCaseAwChartResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTaskCaseAwChartResponse struct{}"
	}

	return strings.Join([]string{"ShowTaskCaseAwChartResponse", string(data)}, " ")
}
