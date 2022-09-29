package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowMessagesResponse struct {

	// 消息列表。
	Messages *[]ShowMessagesRespMessages `json:"messages,omitempty"`

	// 消息总数。
	MessagesCount *int32 `json:"messages_count,omitempty"`

	// 总页数。
	OffsetsCount *int32 `json:"offsets_count,omitempty"`

	// 当前页数。
	Offset         *int32 `json:"offset,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowMessagesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMessagesResponse struct{}"
	}

	return strings.Join([]string{"ShowMessagesResponse", string(data)}, " ")
}
