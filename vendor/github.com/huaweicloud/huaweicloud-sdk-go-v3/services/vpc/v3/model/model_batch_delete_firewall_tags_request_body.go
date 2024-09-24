package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteFirewallTagsRequestBody This is a auto create Body Object
type BatchDeleteFirewallTagsRequestBody struct {

	// 标签列表
	Tags *[]DeleteResourceTagRequestBody `json:"tags,omitempty"`
}

func (o BatchDeleteFirewallTagsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteFirewallTagsRequestBody struct{}"
	}

	return strings.Join([]string{"BatchDeleteFirewallTagsRequestBody", string(data)}, " ")
}
