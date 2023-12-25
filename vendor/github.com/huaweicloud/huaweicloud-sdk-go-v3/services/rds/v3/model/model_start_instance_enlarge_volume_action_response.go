package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartInstanceEnlargeVolumeActionResponse Response Object
type StartInstanceEnlargeVolumeActionResponse struct {

	// 扩容数据库磁盘空间的任务id。 仅磁盘扩容按需实例时会返回该参数。
	JobId *string `json:"job_id,omitempty"`

	// 订单号，磁盘扩容包年包月时返回该参数。
	OrderId        *string `json:"order_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StartInstanceEnlargeVolumeActionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartInstanceEnlargeVolumeActionResponse struct{}"
	}

	return strings.Join([]string{"StartInstanceEnlargeVolumeActionResponse", string(data)}, " ")
}
