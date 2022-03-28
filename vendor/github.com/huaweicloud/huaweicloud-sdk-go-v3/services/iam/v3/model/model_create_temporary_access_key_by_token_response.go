package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTemporaryAccessKeyByTokenResponse struct {
	Credential     *Credential `json:"credential,omitempty"`
	HttpStatusCode int         `json:"-"`
}

func (o CreateTemporaryAccessKeyByTokenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTemporaryAccessKeyByTokenResponse struct{}"
	}

	return strings.Join([]string{"CreateTemporaryAccessKeyByTokenResponse", string(data)}, " ")
}
