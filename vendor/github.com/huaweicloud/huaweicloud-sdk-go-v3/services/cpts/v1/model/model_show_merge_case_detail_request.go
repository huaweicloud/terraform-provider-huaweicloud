package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowMergeCaseDetailRequest Request Object
type ShowMergeCaseDetailRequest struct {

	// 任务运行id（报告id）
	TaskRunId int32 `json:"task_run_id"`

	// 用例运行id
	CaseRunId int32 `json:"case_run_id"`
}

func (o ShowMergeCaseDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMergeCaseDetailRequest struct{}"
	}

	return strings.Join([]string{"ShowMergeCaseDetailRequest", string(data)}, " ")
}
