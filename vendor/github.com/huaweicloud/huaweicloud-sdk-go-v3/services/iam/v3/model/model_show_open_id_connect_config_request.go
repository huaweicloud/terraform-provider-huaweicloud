package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowOpenIdConnectConfigRequest struct {

	// 身份提供商ID
	IdpId string `json:"idp_id"`
}

func (o ShowOpenIdConnectConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowOpenIdConnectConfigRequest struct{}"
	}

	return strings.Join([]string{"ShowOpenIdConnectConfigRequest", string(data)}, " ")
}
