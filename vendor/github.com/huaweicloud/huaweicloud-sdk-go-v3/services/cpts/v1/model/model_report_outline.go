package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ReportOutline struct {
	CoreIndex *CoreIndex `json:"core_index,omitempty"`

	ExceptionResponseSum *ExceptionResponseSum `json:"exception_response_sum,omitempty"`

	NetworkTraffic *NetworkTraffic `json:"network_traffic,omitempty"`

	ResponseCodeSum *ResponseCodeSum `json:"response_code_sum,omitempty"`

	TaskBasicAttribute *TaskBasicAttribute `json:"task_basic_attribute,omitempty"`

	TaskBasicExecutionData *TaskBasicExecutionData `json:"task_basic_execution_data,omitempty"`

	// 响应码详细信息
	ResponseCodeDetails *[]interface{} `json:"response_code_details,omitempty"`

	// SLA数据
	SlaStatistic *interface{} `json:"sla_statistic,omitempty"`

	// 流媒体相关数据
	StreamingMedia *interface{} `json:"streaming_media,omitempty"`
}

func (o ReportOutline) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReportOutline struct{}"
	}

	return strings.Join([]string{"ReportOutline", string(data)}, " ")
}
