package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneRemoveDomainPermissionFromGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o KeystoneRemoveDomainPermissionFromGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneRemoveDomainPermissionFromGroupResponse struct{}"
	}

	return strings.Join([]string{"KeystoneRemoveDomainPermissionFromGroupResponse", string(data)}, " ")
}
