package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneCreateMappingRequestBody struct {
	Mapping *MappingOption `json:"mapping"`
}

func (o KeystoneCreateMappingRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateMappingRequestBody struct{}"
	}

	return strings.Join([]string{"KeystoneCreateMappingRequestBody", string(data)}, " ")
}
