package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartResizeFlavorActionResponse Response Object
type StartResizeFlavorActionResponse struct {

	// 规格变更的任务id。 仅规格变更按需实例时会返回该参数。
	JobId *string `json:"job_id,omitempty"`

	// 订单号，规格变更包年包月时返回该参数。
	OrderId        *string `json:"order_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StartResizeFlavorActionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartResizeFlavorActionResponse struct{}"
	}

	return strings.Join([]string{"StartResizeFlavorActionResponse", string(data)}, " ")
}
