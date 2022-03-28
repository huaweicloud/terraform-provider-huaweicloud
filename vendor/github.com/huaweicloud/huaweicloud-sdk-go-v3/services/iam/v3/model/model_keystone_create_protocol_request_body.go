package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneCreateProtocolRequestBody struct {
	Protocol *ProtocolOption `json:"protocol"`
}

func (o KeystoneCreateProtocolRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateProtocolRequestBody struct{}"
	}

	return strings.Join([]string{"KeystoneCreateProtocolRequestBody", string(data)}, " ")
}
