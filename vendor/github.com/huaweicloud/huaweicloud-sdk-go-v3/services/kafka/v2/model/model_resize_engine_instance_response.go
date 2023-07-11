package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResizeEngineInstanceResponse Response Object
type ResizeEngineInstanceResponse struct {

	// 规格变更任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ResizeEngineInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResizeEngineInstanceResponse struct{}"
	}

	return strings.Join([]string{"ResizeEngineInstanceResponse", string(data)}, " ")
}
