package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListFirewallsByTagsRequestBody This is a auto create Body Object
type ListFirewallsByTagsRequestBody struct {

	// 功能说明：搜索字段，key为要匹配的字段，value为匹配的值 约束：key为固定字典值，不能包含重复的key或不支持的key，当前key仅支持resource_name
	Matches *[]Match `json:"matches,omitempty"`

	// 包含标签，最多包含50个key，每个key下面的value最多10个，每个key对应的value可以为空数组但结构体不能缺失 。 Key不能重复，同一个key中values不能重复。 结果返回包含所有标签的资源列表，key之间是与的关系，key-value结构中value是或的关系。无tag过滤条件时返回全量数据。
	Tags *[]ListTag `json:"tags,omitempty"`
}

func (o ListFirewallsByTagsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListFirewallsByTagsRequestBody struct{}"
	}

	return strings.Join([]string{"ListFirewallsByTagsRequestBody", string(data)}, " ")
}
