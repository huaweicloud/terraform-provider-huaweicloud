package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateFirewallTagsRequestBody This is a auto create Body Object
type BatchCreateFirewallTagsRequestBody struct {

	// 标签列表
	Tags *[]ResourceTag `json:"tags,omitempty"`
}

func (o BatchCreateFirewallTagsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateFirewallTagsRequestBody struct{}"
	}

	return strings.Join([]string{"BatchCreateFirewallTagsRequestBody", string(data)}, " ")
}
