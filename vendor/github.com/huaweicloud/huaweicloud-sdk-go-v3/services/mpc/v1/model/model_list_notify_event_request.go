package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListNotifyEventRequest Request Object
type ListNotifyEventRequest struct {
}

func (o ListNotifyEventRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListNotifyEventRequest struct{}"
	}

	return strings.Join([]string{"ListNotifyEventRequest", string(data)}, " ")
}
