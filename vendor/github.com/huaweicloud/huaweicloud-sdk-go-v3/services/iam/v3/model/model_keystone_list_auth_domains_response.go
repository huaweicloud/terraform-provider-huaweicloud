package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneListAuthDomainsResponse struct {

	// 账号信息列表。
	Domains *[]Domains `json:"domains,omitempty"`

	Links          *LinksSelf `json:"links,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o KeystoneListAuthDomainsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListAuthDomainsResponse struct{}"
	}

	return strings.Join([]string{"KeystoneListAuthDomainsResponse", string(data)}, " ")
}
