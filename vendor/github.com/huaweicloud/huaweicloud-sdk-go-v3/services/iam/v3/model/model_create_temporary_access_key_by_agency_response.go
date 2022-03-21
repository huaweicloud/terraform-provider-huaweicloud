package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTemporaryAccessKeyByAgencyResponse struct {
	Credential     *Credential `json:"credential,omitempty"`
	HttpStatusCode int         `json:"-"`
}

func (o CreateTemporaryAccessKeyByAgencyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTemporaryAccessKeyByAgencyResponse struct{}"
	}

	return strings.Join([]string{"CreateTemporaryAccessKeyByAgencyResponse", string(data)}, " ")
}
