package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneAssociateGroupWithDomainPermissionResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o KeystoneAssociateGroupWithDomainPermissionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneAssociateGroupWithDomainPermissionResponse struct{}"
	}

	return strings.Join([]string{"KeystoneAssociateGroupWithDomainPermissionResponse", string(data)}, " ")
}
