package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartInstanceSingleToHaActionResponse Response Object
type StartInstanceSingleToHaActionResponse struct {

	// 单机转主备的任务id。 仅按需实例单机转主备时会返回该参数。
	JobId *string `json:"job_id,omitempty"`

	// 订单号，包年包月单机转主备时返回该参数。
	OrderId        *string `json:"order_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StartInstanceSingleToHaActionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartInstanceSingleToHaActionResponse struct{}"
	}

	return strings.Join([]string{"StartInstanceSingleToHaActionResponse", string(data)}, " ")
}
