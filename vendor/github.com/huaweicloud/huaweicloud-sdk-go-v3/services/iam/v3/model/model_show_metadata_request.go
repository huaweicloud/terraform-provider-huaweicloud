package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowMetadataRequest struct {

	// 身份提供商ID。
	IdpId string `json:"idp_id"`

	// 协议ID。
	ProtocolId string `json:"protocol_id"`
}

func (o ShowMetadataRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMetadataRequest struct{}"
	}

	return strings.Join([]string{"ShowMetadataRequest", string(data)}, " ")
}
