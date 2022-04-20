package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TaskRunInfo
type TaskRunInfo struct {
	// id

	Id *int32 `json:"id,omitempty"`
	// run_type

	RunType *int32 `json:"run_type,omitempty"`
}

func (o TaskRunInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskRunInfo struct{}"
	}

	return strings.Join([]string{"TaskRunInfo", string(data)}, " ")
}
