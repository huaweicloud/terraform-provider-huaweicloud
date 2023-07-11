package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopBatchTaskResponse Response Object
type StopBatchTaskResponse struct {

	// 批量操作目标结果集合
	Targets        *[]BatchTargetResult `json:"targets,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o StopBatchTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopBatchTaskResponse struct{}"
	}

	return strings.Join([]string{"StopBatchTaskResponse", string(data)}, " ")
}
