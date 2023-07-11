package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateReassignmentTaskResponse Response Object
type CreateReassignmentTaskResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateReassignmentTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateReassignmentTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateReassignmentTaskResponse", string(data)}, " ")
}
