package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneCreateMappingResponse struct {
	Mapping        *MappingResult `json:"mapping,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o KeystoneCreateMappingResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateMappingResponse struct{}"
	}

	return strings.Join([]string{"KeystoneCreateMappingResponse", string(data)}, " ")
}
