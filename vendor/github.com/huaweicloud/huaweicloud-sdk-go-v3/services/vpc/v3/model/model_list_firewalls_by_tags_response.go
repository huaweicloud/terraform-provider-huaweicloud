package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListFirewallsByTagsResponse Response Object
type ListFirewallsByTagsResponse struct {

	// 资源列表
	Resources *[]ListResourceResp `json:"resources,omitempty"`

	// 资源数量
	TotalCount *int32 `json:"total_count,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListFirewallsByTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListFirewallsByTagsResponse struct{}"
	}

	return strings.Join([]string{"ListFirewallsByTagsResponse", string(data)}, " ")
}
