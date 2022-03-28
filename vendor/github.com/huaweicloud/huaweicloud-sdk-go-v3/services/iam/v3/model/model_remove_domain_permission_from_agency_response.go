package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type RemoveDomainPermissionFromAgencyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RemoveDomainPermissionFromAgencyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveDomainPermissionFromAgencyResponse struct{}"
	}

	return strings.Join([]string{"RemoveDomainPermissionFromAgencyResponse", string(data)}, " ")
}
