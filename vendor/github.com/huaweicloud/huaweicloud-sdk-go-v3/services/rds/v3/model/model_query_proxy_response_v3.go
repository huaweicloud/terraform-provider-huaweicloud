package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// QueryProxyResponseV3 查询数据库代理详情信息返回体。
type QueryProxyResponseV3 struct {
	Proxy *ProxyInfo `json:"proxy,omitempty"`

	MasterInstance *InstanceInfo `json:"master_instance,omitempty"`

	// 数据库只读实例信息。
	ReadonlyInstances *[]InstanceInfo `json:"readonly_instances,omitempty"`

	// 安全组是否放通该数据库代理到数据库的网络地址。
	ProxySecurityGroupCheckResult *bool `json:"proxy_security_group_check_result,omitempty"`
}

func (o QueryProxyResponseV3) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QueryProxyResponseV3 struct{}"
	}

	return strings.Join([]string{"QueryProxyResponseV3", string(data)}, " ")
}
