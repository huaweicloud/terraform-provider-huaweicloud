package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExportPrivateKeyRequest Request Object
type ExportPrivateKeyRequest struct {
	Body *ExportPrivateKeyRequestBody `json:"body,omitempty"`
}

func (o ExportPrivateKeyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExportPrivateKeyRequest struct{}"
	}

	return strings.Join([]string{"ExportPrivateKeyRequest", string(data)}, " ")
}
