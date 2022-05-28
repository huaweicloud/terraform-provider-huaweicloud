package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowPreheatingAssetRequest struct {

	// 任务ID
	TaskId string `json:"task_id"`
}

func (o ShowPreheatingAssetRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPreheatingAssetRequest struct{}"
	}

	return strings.Join([]string{"ShowPreheatingAssetRequest", string(data)}, " ")
}
