package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateKeypairResponse struct {
	Keypair        *CreateKeypairResp `json:"keypair,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o CreateKeypairResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateKeypairResponse struct{}"
	}

	return strings.Join([]string{"CreateKeypairResponse", string(data)}, " ")
}
