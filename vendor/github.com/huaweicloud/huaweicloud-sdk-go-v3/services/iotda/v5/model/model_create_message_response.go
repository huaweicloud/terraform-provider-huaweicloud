package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateMessageResponse struct {

	// 消息id，由用户生成（推荐使用UUID），同一个设备下唯一， 如果用户不填写，则由物联网平台生成。
	MessageId *string `json:"message_id,omitempty"`

	Result         *MessageResult `json:"result,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o CreateMessageResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMessageResponse struct{}"
	}

	return strings.Join([]string{"CreateMessageResponse", string(data)}, " ")
}
