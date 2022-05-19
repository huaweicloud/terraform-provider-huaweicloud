package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListEditingJobResponse struct {

	// 任务总数
	Total *int32 `json:"total,omitempty"`

	// 任务列表
	Jobs           *[]EditingJob `json:"jobs,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o ListEditingJobResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEditingJobResponse struct{}"
	}

	return strings.Join([]string{"ListEditingJobResponse", string(data)}, " ")
}
