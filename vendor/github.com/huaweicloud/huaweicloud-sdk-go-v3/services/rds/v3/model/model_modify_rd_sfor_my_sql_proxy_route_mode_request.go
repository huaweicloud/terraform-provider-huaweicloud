package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyRdSforMySqlProxyRouteModeRequest Request Object
type ModifyRdSforMySqlProxyRouteModeRequest struct {
	ContentType *string `json:"Content-Type,omitempty"`

	// 实例ID，严格匹配UUID规则。
	InstanceId string `json:"instance_id"`

	// 数据库代理ID，严格匹配UUID规则。
	ProxyId string `json:"proxy_id"`

	// 语言。
	XLanguage *string `json:"X-Language,omitempty"`

	Body *ModifyMySqlProxyRouteModeRequest `json:"body,omitempty"`
}

func (o ModifyRdSforMySqlProxyRouteModeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyRdSforMySqlProxyRouteModeRequest struct{}"
	}

	return strings.Join([]string{"ModifyRdSforMySqlProxyRouteModeRequest", string(data)}, " ")
}
