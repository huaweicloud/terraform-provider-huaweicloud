package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDomainKeyChainResponse Response Object
type UpdateDomainKeyChainResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateDomainKeyChainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainKeyChainResponse struct{}"
	}

	return strings.Join([]string{"UpdateDomainKeyChainResponse", string(data)}, " ")
}
