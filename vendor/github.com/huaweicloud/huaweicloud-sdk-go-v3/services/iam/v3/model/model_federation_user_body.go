package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 联邦用户信息。
type FederationUserBody struct {
	OsFederation *OsFederationInfo `json:"OS-FEDERATION"`

	Domain *DomainInfo `json:"domain"`

	// user id。
	Id *string `json:"id,omitempty"`

	// user name。
	Name *string `json:"name,omitempty"`
}

func (o FederationUserBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FederationUserBody struct{}"
	}

	return strings.Join([]string{"FederationUserBody", string(data)}, " ")
}
