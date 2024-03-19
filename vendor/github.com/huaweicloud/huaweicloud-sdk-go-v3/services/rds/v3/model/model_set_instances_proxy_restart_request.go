package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetInstancesProxyRestartRequest Request Object
type SetInstancesProxyRestartRequest struct {
	ContentType *string `json:"Content-Type,omitempty"`

	// 实例ID，严格匹配UUID规则。
	InstanceId string `json:"instance_id"`

	// 数据库代理ID，严格匹配UUID规则。
	ProxyId string `json:"proxy_id"`

	// 语言。
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o SetInstancesProxyRestartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetInstancesProxyRestartRequest struct{}"
	}

	return strings.Join([]string{"SetInstancesProxyRestartRequest", string(data)}, " ")
}
