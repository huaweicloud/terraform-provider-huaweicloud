package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CountFirewallsByTagsResponse Response Object
type CountFirewallsByTagsResponse struct {

	// 请求ID
	RequestId *string `json:"request_id,omitempty"`

	// 资源数量
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CountFirewallsByTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CountFirewallsByTagsResponse struct{}"
	}

	return strings.Join([]string{"CountFirewallsByTagsResponse", string(data)}, " ")
}
