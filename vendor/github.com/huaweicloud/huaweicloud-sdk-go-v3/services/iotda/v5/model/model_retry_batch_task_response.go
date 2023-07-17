package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RetryBatchTaskResponse Response Object
type RetryBatchTaskResponse struct {

	// 批量操作目标结果集合
	Targets        *[]BatchTargetResult `json:"targets,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o RetryBatchTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RetryBatchTaskResponse struct{}"
	}

	return strings.Join([]string{"RetryBatchTaskResponse", string(data)}, " ")
}
