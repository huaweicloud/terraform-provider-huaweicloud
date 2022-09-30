package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowInstanceMessagesResponse struct {

	// 消息列表。
	Messages *[]MessagesEntity `json:"messages,omitempty"`

	// 消息总条数。
	Total *int64 `json:"total,omitempty"`

	// 消息条数。
	Size           *int64 `json:"size,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowInstanceMessagesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceMessagesResponse struct{}"
	}

	return strings.Join([]string{"ShowInstanceMessagesResponse", string(data)}, " ")
}
