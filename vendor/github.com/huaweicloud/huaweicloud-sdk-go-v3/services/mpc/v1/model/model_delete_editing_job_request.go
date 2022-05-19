package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteEditingJobRequest struct {

	// 任务ID
	JobId string `json:"job_id"`
}

func (o DeleteEditingJobRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteEditingJobRequest struct{}"
	}

	return strings.Join([]string{"DeleteEditingJobRequest", string(data)}, " ")
}
