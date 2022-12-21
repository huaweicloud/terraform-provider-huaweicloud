package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListNotifiedHistoriesRequest struct {

	// 告警流水号
	EventSn *string `json:"event_sn,omitempty"`
}

func (o ListNotifiedHistoriesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListNotifiedHistoriesRequest struct{}"
	}

	return strings.Join([]string{"ListNotifiedHistoriesRequest", string(data)}, " ")
}
