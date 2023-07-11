package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListNotifiedHistoriesResponse Response Object
type ListNotifiedHistoriesResponse struct {

	// 告警流水号
	EventSn *string `json:"event_sn,omitempty"`

	// 通知结果
	Notifications  *[]Notifications `json:"notifications,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListNotifiedHistoriesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListNotifiedHistoriesResponse struct{}"
	}

	return strings.Join([]string{"ListNotifiedHistoriesResponse", string(data)}, " ")
}
