package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneUpdateProtocolResponse struct {
	Protocol       *ProtocolResult `json:"protocol,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o KeystoneUpdateProtocolResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneUpdateProtocolResponse struct{}"
	}

	return strings.Join([]string{"KeystoneUpdateProtocolResponse", string(data)}, " ")
}
