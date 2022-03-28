package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreatePermanentAccessKeyResponse struct {
	Credential     *CreateCredentialResult `json:"credential,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o CreatePermanentAccessKeyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePermanentAccessKeyResponse struct{}"
	}

	return strings.Join([]string{"CreatePermanentAccessKeyResponse", string(data)}, " ")
}
