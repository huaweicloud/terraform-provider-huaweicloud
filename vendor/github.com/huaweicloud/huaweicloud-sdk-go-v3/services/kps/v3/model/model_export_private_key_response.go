package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExportPrivateKeyResponse Response Object
type ExportPrivateKeyResponse struct {
	Keypair        *ExportPrivateKeyKeypairBean `json:"keypair,omitempty"`
	HttpStatusCode int                          `json:"-"`
}

func (o ExportPrivateKeyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExportPrivateKeyResponse struct{}"
	}

	return strings.Join([]string{"ExportPrivateKeyResponse", string(data)}, " ")
}
