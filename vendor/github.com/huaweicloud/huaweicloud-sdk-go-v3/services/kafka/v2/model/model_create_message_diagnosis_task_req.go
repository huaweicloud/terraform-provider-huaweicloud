package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateMessageDiagnosisTaskReq struct {

	// 消费组名称
	GroupName string `json:"group_name"`

	// topic名称
	TopicName string `json:"topic_name"`
}

func (o CreateMessageDiagnosisTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMessageDiagnosisTaskReq struct{}"
	}

	return strings.Join([]string{"CreateMessageDiagnosisTaskReq", string(data)}, " ")
}
