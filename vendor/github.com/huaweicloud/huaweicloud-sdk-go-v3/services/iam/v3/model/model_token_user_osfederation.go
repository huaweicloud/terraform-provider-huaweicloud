package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type TokenUserOsfederation struct {

	// 用户组信息列表。
	Groups []OsfederationGroups `json:"groups"`

	IdentityProvider *OsfederationIdentityprovider `json:"identity_provider"`

	Protocol *OsfederationProtocol `json:"protocol"`
}

func (o TokenUserOsfederation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TokenUserOsfederation struct{}"
	}

	return strings.Join([]string{"TokenUserOsfederation", string(data)}, " ")
}
