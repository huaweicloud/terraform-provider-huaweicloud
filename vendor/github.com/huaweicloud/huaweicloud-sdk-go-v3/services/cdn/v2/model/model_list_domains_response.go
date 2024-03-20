package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDomainsResponse Response Object
type ListDomainsResponse struct {

	// 总条数。
	Total *int32 `json:"total,omitempty"`

	// 域名信息。
	Domains *[]Domains `json:"domains,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListDomainsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDomainsResponse struct{}"
	}

	return strings.Join([]string{"ListDomainsResponse", string(data)}, " ")
}
