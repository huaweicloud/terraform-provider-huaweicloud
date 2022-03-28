package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowOpenIdConnectConfigResponse struct {
	OpenidConnectConfig *OpenIdConnectConfig `json:"openid_connect_config,omitempty"`
	HttpStatusCode      int                  `json:"-"`
}

func (o ShowOpenIdConnectConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowOpenIdConnectConfigResponse struct{}"
	}

	return strings.Join([]string{"ShowOpenIdConnectConfigResponse", string(data)}, " ")
}
