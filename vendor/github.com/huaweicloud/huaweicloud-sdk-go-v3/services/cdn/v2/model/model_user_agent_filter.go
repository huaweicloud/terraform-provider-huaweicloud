package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UserAgentFilter UA黑白名单设置。
type UserAgentFilter struct {

	// UA黑白名单类型 off：关闭UA黑白名单; black：UA黑名单; white：UA白名单;
	Type string `json:"type"`

	// 配置UA黑白名单，当type=off时，非必传。最多配置10条规则，单条规则不超过100个字符，多条规则用“,”分割。
	Value *string `json:"value,omitempty"`

	// 是否包含空User-Agent，true:包含，false：不包含。空User-Agent是指没有User-Agent字段或者该字段的值为空。如果黑名单且该字段值为true，则表示空User-Agent不允许访问，如果是白名单且该字段值为true，则表示空User-Agent允许访问。设置User-Agent黑名单时，默认值为false，设置User-Agent白名单时，默认值为true。
	IncludeEmpty *bool `json:"include_empty,omitempty"`

	// 配置UA黑白名单，当type=off时，非必传。最多配置10条规则，单条规则不超过100个字符,同时配置value和ua_list时，ua_list生效。
	UaList *[]string `json:"ua_list,omitempty"`
}

func (o UserAgentFilter) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserAgentFilter struct{}"
	}

	return strings.Join([]string{"UserAgentFilter", string(data)}, " ")
}
