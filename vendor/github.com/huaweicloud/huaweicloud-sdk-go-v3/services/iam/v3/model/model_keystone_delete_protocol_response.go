package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneDeleteProtocolResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o KeystoneDeleteProtocolResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneDeleteProtocolResponse struct{}"
	}

	return strings.Join([]string{"KeystoneDeleteProtocolResponse", string(data)}, " ")
}
