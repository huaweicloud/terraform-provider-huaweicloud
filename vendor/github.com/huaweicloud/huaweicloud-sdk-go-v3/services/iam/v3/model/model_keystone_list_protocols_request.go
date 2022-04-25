package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListProtocolsRequest struct {

	// 身份提供商ID。
	IdpId string `json:"idp_id"`
}

func (o KeystoneListProtocolsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListProtocolsRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListProtocolsRequest", string(data)}, " ")
}
