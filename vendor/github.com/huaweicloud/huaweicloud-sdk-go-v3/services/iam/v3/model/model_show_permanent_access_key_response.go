package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowPermanentAccessKeyResponse struct {
	Credential     *ShowCredential `json:"credential,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ShowPermanentAccessKeyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPermanentAccessKeyResponse struct{}"
	}

	return strings.Join([]string{"ShowPermanentAccessKeyResponse", string(data)}, " ")
}
