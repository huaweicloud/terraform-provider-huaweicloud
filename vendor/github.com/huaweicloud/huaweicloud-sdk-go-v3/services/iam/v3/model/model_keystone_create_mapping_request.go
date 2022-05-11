package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneCreateMappingRequest struct {

	// 待注册的映射ID。
	Id string `json:"id"`

	Body *KeystoneCreateMappingRequestBody `json:"body,omitempty"`
}

func (o KeystoneCreateMappingRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateMappingRequest struct{}"
	}

	return strings.Join([]string{"KeystoneCreateMappingRequest", string(data)}, " ")
}
