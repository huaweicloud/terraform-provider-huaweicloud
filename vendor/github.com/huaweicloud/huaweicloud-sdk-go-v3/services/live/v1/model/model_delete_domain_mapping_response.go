package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type DeleteDomainMappingResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteDomainMappingResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDomainMappingResponse struct{}"
	}

	return strings.Join([]string{"DeleteDomainMappingResponse", string(data)}, " ")
}
