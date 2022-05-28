package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowTakeOverTaskDetailsResponse struct {

	// 总数。
	Total *int32 `json:"total,omitempty"`

	// 任务ID。
	TaskId *string `json:"task_id,omitempty"`

	// 任务状态。
	TaskStatus *string `json:"task_status,omitempty"`

	// 媒资信息。
	Assets         *[]AssetDetails `json:"assets,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ShowTakeOverTaskDetailsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTakeOverTaskDetailsResponse struct{}"
	}

	return strings.Join([]string{"ShowTakeOverTaskDetailsResponse", string(data)}, " ")
}
