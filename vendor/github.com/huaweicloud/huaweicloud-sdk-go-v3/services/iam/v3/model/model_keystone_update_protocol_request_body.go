package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneUpdateProtocolRequestBody struct {
	Protocol *ProtocolOption `json:"protocol"`
}

func (o KeystoneUpdateProtocolRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneUpdateProtocolRequestBody struct{}"
	}

	return strings.Join([]string{"KeystoneUpdateProtocolRequestBody", string(data)}, " ")
}
