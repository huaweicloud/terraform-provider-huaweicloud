package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ImportPrivateKeyResponse Response Object
type ImportPrivateKeyResponse struct {
	Keypair        *ImportPrivateKeyAction `json:"keypair,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o ImportPrivateKeyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImportPrivateKeyResponse struct{}"
	}

	return strings.Join([]string{"ImportPrivateKeyResponse", string(data)}, " ")
}
