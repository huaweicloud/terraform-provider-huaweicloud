package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListAuthDomainsRequest struct {
}

func (o KeystoneListAuthDomainsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListAuthDomainsRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListAuthDomainsRequest", string(data)}, " ")
}
