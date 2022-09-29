package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ResizeInstanceResponse struct {

	// 规格变更任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ResizeInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResizeInstanceResponse struct{}"
	}

	return strings.Join([]string{"ResizeInstanceResponse", string(data)}, " ")
}
