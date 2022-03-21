package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type CreatePermanentAccessKeyRequestBody struct {
	Credential *CreateCredentialOption `json:"credential"`
}

func (o CreatePermanentAccessKeyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePermanentAccessKeyRequestBody struct{}"
	}

	return strings.Join([]string{"CreatePermanentAccessKeyRequestBody", string(data)}, " ")
}
