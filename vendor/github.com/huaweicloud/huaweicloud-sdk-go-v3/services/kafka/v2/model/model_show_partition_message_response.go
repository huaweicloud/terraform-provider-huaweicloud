package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPartitionMessageResponse Response Object
type ShowPartitionMessageResponse struct {

	// 消息列表。
	Message        *[]ShowPartitionMessageEntity `json:"message,omitempty"`
	HttpStatusCode int                           `json:"-"`
}

func (o ShowPartitionMessageResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPartitionMessageResponse struct{}"
	}

	return strings.Join([]string{"ShowPartitionMessageResponse", string(data)}, " ")
}
