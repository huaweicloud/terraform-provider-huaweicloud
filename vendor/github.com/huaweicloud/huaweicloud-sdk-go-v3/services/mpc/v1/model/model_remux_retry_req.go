package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RemuxRetryReq struct {

	// 任务Id。
	TaskId *string `json:"task_id,omitempty"`
}

func (o RemuxRetryReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemuxRetryReq struct{}"
	}

	return strings.Join([]string{"RemuxRetryReq", string(data)}, " ")
}
