package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneCreateGroupRequestBody struct {
	Group *KeystoneCreateGroupOption `json:"group"`
}

func (o KeystoneCreateGroupRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateGroupRequestBody struct{}"
	}

	return strings.Join([]string{"KeystoneCreateGroupRequestBody", string(data)}, " ")
}
