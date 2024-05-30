package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResetBridgeSecretResponse Response Object
type ResetBridgeSecretResponse struct {

	// 网桥ID
	BridgeId *string `json:"bridge_id,omitempty"`

	// 网桥密钥。
	Secret         *string `json:"secret,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ResetBridgeSecretResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetBridgeSecretResponse struct{}"
	}

	return strings.Join([]string{"ResetBridgeSecretResponse", string(data)}, " ")
}
