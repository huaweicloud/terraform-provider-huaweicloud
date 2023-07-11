package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDomainKeyChainResponse Response Object
type DeleteDomainKeyChainResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteDomainKeyChainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDomainKeyChainResponse struct{}"
	}

	return strings.Join([]string{"DeleteDomainKeyChainResponse", string(data)}, " ")
}
