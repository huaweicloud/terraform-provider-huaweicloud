package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DisableDomainResponse Response Object
type DisableDomainResponse struct {
	Domain *DomainsWithPort `json:"domain,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DisableDomainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DisableDomainResponse struct{}"
	}

	return strings.Join([]string{"DisableDomainResponse", string(data)}, " ")
}
