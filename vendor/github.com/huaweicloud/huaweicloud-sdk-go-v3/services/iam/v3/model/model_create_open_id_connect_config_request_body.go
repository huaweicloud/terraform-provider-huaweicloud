package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 请求体
type CreateOpenIdConnectConfigRequestBody struct {
	OpenidConnectConfig *CreateOpenIdConnectConfig `json:"openid_connect_config"`
}

func (o CreateOpenIdConnectConfigRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateOpenIdConnectConfigRequestBody struct{}"
	}

	return strings.Join([]string{"CreateOpenIdConnectConfigRequestBody", string(data)}, " ")
}
