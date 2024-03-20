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
