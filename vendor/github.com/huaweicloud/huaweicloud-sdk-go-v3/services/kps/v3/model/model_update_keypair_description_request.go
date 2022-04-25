package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateKeypairDescriptionRequest struct {

	// 密钥对名称
	KeypairName string `json:"keypair_name"`

	Body *UpdateKeypairDescriptionRequestBody `json:"body,omitempty"`
}

func (o UpdateKeypairDescriptionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateKeypairDescriptionRequest struct{}"
	}

	return strings.Join([]string{"UpdateKeypairDescriptionRequest", string(data)}, " ")
}
