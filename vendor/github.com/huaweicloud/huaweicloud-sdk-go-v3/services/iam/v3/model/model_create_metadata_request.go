package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateMetadataRequest struct {

	// 身份提供商ID。
	IdpId string `json:"idp_id"`

	// 协议ID。
	ProtocolId string `json:"protocol_id"`

	Body *CreateMetadataRequestBody `json:"body,omitempty"`
}

func (o CreateMetadataRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMetadataRequest struct{}"
	}

	return strings.Join([]string{"CreateMetadataRequest", string(data)}, " ")
}
