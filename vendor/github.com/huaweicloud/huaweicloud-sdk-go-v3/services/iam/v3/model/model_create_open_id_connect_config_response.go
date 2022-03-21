package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateOpenIdConnectConfigResponse struct {
	OpenidConnectConfig *CreateOpenIdConnectConfig `json:"openid_connect_config,omitempty"`
	HttpStatusCode      int                        `json:"-"`
}

func (o CreateOpenIdConnectConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateOpenIdConnectConfigResponse struct{}"
	}

	return strings.Join([]string{"CreateOpenIdConnectConfigResponse", string(data)}, " ")
}
