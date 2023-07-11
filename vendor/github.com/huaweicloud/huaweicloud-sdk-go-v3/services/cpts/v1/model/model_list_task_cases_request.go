package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTaskCasesRequest Request Object
type ListTaskCasesRequest struct {

	// 工程id
	TestSuitId int32 `json:"test_suit_id"`

	// 任务id
	TaskId int32 `json:"task_id"`
}

func (o ListTaskCasesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTaskCasesRequest struct{}"
	}

	return strings.Join([]string{"ListTaskCasesRequest", string(data)}, " ")
}
