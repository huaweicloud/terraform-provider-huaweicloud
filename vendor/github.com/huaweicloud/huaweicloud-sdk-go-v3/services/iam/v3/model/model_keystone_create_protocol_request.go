package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneCreateProtocolRequest struct {

	// 身份提供商ID。
	IdpId string `json:"idp_id"`

	// 待注册的协议ID。
	ProtocolId string `json:"protocol_id"`

	Body *KeystoneCreateProtocolRequestBody `json:"body,omitempty"`
}

func (o KeystoneCreateProtocolRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateProtocolRequest struct{}"
	}

	return strings.Join([]string{"KeystoneCreateProtocolRequest", string(data)}, " ")
}
