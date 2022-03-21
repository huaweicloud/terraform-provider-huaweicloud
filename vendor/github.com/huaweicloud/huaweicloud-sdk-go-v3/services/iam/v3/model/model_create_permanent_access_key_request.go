package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreatePermanentAccessKeyRequest struct {
	Body *CreatePermanentAccessKeyRequestBody `json:"body,omitempty"`
}

func (o CreatePermanentAccessKeyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePermanentAccessKeyRequest struct{}"
	}

	return strings.Join([]string{"CreatePermanentAccessKeyRequest", string(data)}, " ")
}
