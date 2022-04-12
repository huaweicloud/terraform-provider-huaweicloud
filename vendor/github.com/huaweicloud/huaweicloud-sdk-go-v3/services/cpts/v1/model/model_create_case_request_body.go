package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateCaseRequestBody
type CreateCaseRequestBody struct {
	// name

	Name string `json:"name"`
	// type (0-常规用例,1-视频流用例,2-预制用例)

	Type int32 `json:"type"`
	// task_id

	TaskId int32 `json:"task_id"`
}

func (o CreateCaseRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCaseRequestBody struct{}"
	}

	return strings.Join([]string{"CreateCaseRequestBody", string(data)}, " ")
}
