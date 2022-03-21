package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListFederationDomainsRequest struct {
}

func (o KeystoneListFederationDomainsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListFederationDomainsRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListFederationDomainsRequest", string(data)}, " ")
}
