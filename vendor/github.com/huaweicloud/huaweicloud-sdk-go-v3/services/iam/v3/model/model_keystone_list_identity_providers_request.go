package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// KeystoneListIdentityProvidersRequest Request Object
type KeystoneListIdentityProvidersRequest struct {
}

func (o KeystoneListIdentityProvidersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListIdentityProvidersRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListIdentityProvidersRequest", string(data)}, " ")
}
