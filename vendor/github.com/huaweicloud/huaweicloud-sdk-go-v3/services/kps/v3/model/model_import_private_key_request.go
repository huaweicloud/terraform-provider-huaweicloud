package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ImportPrivateKeyRequest Request Object
type ImportPrivateKeyRequest struct {
	Body *ImportPrivateKeyRequestBody `json:"body,omitempty"`
}

func (o ImportPrivateKeyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImportPrivateKeyRequest struct{}"
	}

	return strings.Join([]string{"ImportPrivateKeyRequest", string(data)}, " ")
}
