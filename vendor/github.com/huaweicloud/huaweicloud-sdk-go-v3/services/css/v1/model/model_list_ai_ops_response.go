package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAiOpsResponse Response Object
type ListAiOpsResponse struct {

	// 检测任务个数。
	TotalSize *int32 `json:"total_size,omitempty"`

	// 创建一个集群检测任务。
	AiopsList      *[]ListAiOpsRequestBodyAiopsList `json:"aiops_list,omitempty"`
	HttpStatusCode int                              `json:"-"`
}

func (o ListAiOpsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAiOpsResponse struct{}"
	}

	return strings.Join([]string{"ListAiOpsResponse", string(data)}, " ")
}
