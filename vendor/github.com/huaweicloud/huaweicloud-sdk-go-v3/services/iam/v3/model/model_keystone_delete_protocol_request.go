package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneDeleteProtocolRequest struct {

	// 身份提供商ID。
	IdpId string `json:"idp_id"`

	// 待删除的协议ID。
	ProtocolId string `json:"protocol_id"`
}

func (o KeystoneDeleteProtocolRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneDeleteProtocolRequest struct{}"
	}

	return strings.Join([]string{"KeystoneDeleteProtocolRequest", string(data)}, " ")
}
