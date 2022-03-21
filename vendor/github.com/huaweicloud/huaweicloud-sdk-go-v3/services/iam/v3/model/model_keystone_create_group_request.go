package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneCreateGroupRequest struct {
	Body *KeystoneCreateGroupRequestBody `json:"body,omitempty"`
}

func (o KeystoneCreateGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateGroupRequest struct{}"
	}

	return strings.Join([]string{"KeystoneCreateGroupRequest", string(data)}, " ")
}
