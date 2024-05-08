package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListMessageDiagnosisReportsRequest Request Object
type ListMessageDiagnosisReportsRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 偏移量，表示查询该偏移量后面的记录
	Offset *int32 `json:"offset,omitempty"`

	// 查询返回记录的数量限制
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListMessageDiagnosisReportsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMessageDiagnosisReportsRequest struct{}"
	}

	return strings.Join([]string{"ListMessageDiagnosisReportsRequest", string(data)}, " ")
}
