package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateOpenIdConnectConfigRequestBody struct {
	OpenidConnectConfig *UpdateOpenIdConnectConfig `json:"openid_connect_config"`
}

func (o UpdateOpenIdConnectConfigRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateOpenIdConnectConfigRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateOpenIdConnectConfigRequestBody", string(data)}, " ")
}
