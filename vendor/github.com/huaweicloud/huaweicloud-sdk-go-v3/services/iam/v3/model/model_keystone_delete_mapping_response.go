package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneDeleteMappingResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o KeystoneDeleteMappingResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneDeleteMappingResponse struct{}"
	}

	return strings.Join([]string{"KeystoneDeleteMappingResponse", string(data)}, " ")
}
