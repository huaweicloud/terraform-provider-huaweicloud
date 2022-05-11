package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneListFederationDomainsResponse struct {

	// 账号信息列表。
	Domains *[]Domains `json:"domains,omitempty"`

	Links          *LinksSelf `json:"links,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o KeystoneListFederationDomainsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListFederationDomainsResponse struct{}"
	}

	return strings.Join([]string{"KeystoneListFederationDomainsResponse", string(data)}, " ")
}
