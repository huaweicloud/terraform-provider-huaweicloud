package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDomainResponse Response Object
type DeleteDomainResponse struct {
	Domain *DomainsWithPort `json:"domain,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteDomainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDomainResponse struct{}"
	}

	return strings.Join([]string{"DeleteDomainResponse", string(data)}, " ")
}
