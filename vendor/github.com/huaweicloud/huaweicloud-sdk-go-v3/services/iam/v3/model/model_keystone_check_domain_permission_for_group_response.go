package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneCheckDomainPermissionForGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o KeystoneCheckDomainPermissionForGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCheckDomainPermissionForGroupResponse struct{}"
	}

	return strings.Join([]string{"KeystoneCheckDomainPermissionForGroupResponse", string(data)}, " ")
}
