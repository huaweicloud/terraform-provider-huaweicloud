package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListFirewallTagsResponse Response Object
type ListFirewallTagsResponse struct {

	// tag对象列表
	Tags *[]ListTag `json:"tags,omitempty"`

	// 请求ID
	RequestId *string `json:"request_id,omitempty"`

	// 资源数量
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListFirewallTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListFirewallTagsResponse struct{}"
	}

	return strings.Join([]string{"ListFirewallTagsResponse", string(data)}, " ")
}
