package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ProxyInfo 数据库代理信息。
type ProxyInfo struct {

	// 数据库代理实例ID。
	PoolId *string `json:"pool_id,omitempty"`

	// 数据库代理状态。  取值: NORMAL：表示数据库代理正常。 ENABLING：表示数据库代理正在开启。 DISABLING：表示数据库代理正在关闭。 CHANGING_NODE_NUM：表示数据库代理正在调整节点数量。 SCALING: 表示数据库代理正在规格变更。 UPGRADING: 表示数据库代理正在升级内核版本。 IPMODIFYING: 表示数据库代理正在修改读写分离地址。 RESTARTING: 表示数据库代理正在重启进程。 TRANSACTION_SPLITTING: 表示数据库代理正在变更事务拆分功能状态。 CONNECTION_POOL_SWITCH_OPERATING: 表示数据库代理正在变更会话连接池类型。 PORT_MODIFYING: 表示数据库代理正在修改端口。 PROXY_SSL_SWITCHING: 表示数据库代理正在变更SSL状态。 ALT_SWITCH_OPERATING: 表示数据库代理正在变更ALT状态。 CHANGING_RESOURCES: 表示数据库代理正在进行资源变更。 NORMAL: 表示数据库代理正常。 ABNORMAL: 表示数据库代理异常。 FAILED: 表示数据库代理创建失败。 FROZEN: 表示数据库代理已冻结。
	Status *string `json:"status,omitempty"`

	// 读写分离地址。
	Address *string `json:"address,omitempty"`

	// 端口号。
	Port *int32 `json:"port,omitempty"`

	// 延时阈值，单位：秒。
	DelayThresholdInSeconds *int32 `json:"delay_threshold_in_seconds,omitempty"`

	// 数据库代理规格的CPU大小。
	Cpu *string `json:"cpu,omitempty"`

	// 数据库代理规格的内存大小。
	Mem *string `json:"mem,omitempty"`

	// 数据库代理节点个数。
	NodeNum *int32 `json:"node_num,omitempty"`

	// 数据库代理节点信息列表。
	Nodes *[]ProxyInfoNodes `json:"nodes,omitempty"`

	// 数据库代理集群模式。 取值：     Cluster：集群模式     Ha：主备模式
	Mode *string `json:"mode,omitempty"`

	FlavorInfo *ProxyInfoFlavorInfo `json:"flavor_info,omitempty"`

	// 数据库代理事务拆分开关状态。  true：开启。  false：关闭。
	TransactionSplit *string `json:"transaction_split,omitempty"`

	// 连接池类型。  取值范围:  CLOSED: 关闭连接池。  SESSION: 开启会话级连接池。
	ConnectionPoolType *string `json:"connection_pool_type,omitempty"`

	// 数据库代理计费模式。  取值范围： 0:按需计费 1:包周期计费
	PayMode *string `json:"pay_mode,omitempty"`

	// 数据库代理名称。
	Name *string `json:"name,omitempty"`

	// 数据库代理读写模式。 取值范围：     readwrite 读写模式     readonly 只读模式
	ProxyMode *string `json:"proxy_mode,omitempty"`

	// 数据库代理读写分离地址内网域名。 该字段为空表示未申请读写内网域名。
	DnsName *string `json:"dns_name,omitempty"`

	// 数据库代理实例所属子网ID。
	SubnetId *string `json:"subnet_id,omitempty"`

	// 数据库代理秒级监控状态。
	SecondsLevelMonitorFunStatus *ProxyInfoSecondsLevelMonitorFunStatus `json:"seconds_level_monitor_fun_status,omitempty"`

	// ALT开关状态。
	AltFlag *bool `json:"alt_flag,omitempty"`

	// 是否强制读路由到只读。
	ForceReadOnly *bool `json:"force_read_only,omitempty"`

	// 数据库代理路由模式。 取值范围:     0：表示权重负载模式。     1：表示负载均衡模式（数据库主实例不接受读请求）。     2：表示负载均衡模式（数据库主实例接受读请求）。
	RouteMode *int32 `json:"route_mode,omitempty"`

	// ssl开关状态。
	SslOption *bool `json:"ssl_option,omitempty"`

	// 数据库代理是否支持开启负载均衡路由模式。
	SupportBalanceRouteMode *bool `json:"support_balance_route_mode,omitempty"`

	// 数据库代理是否支持开启ssl功能。
	SupportProxySsl *bool `json:"support_proxy_ssl,omitempty"`

	// 数据库代理是否支持切换会话连接池类型。
	SupportSwitchConnectionPoolType *bool `json:"support_switch_connection_pool_type,omitempty"`

	// 数据库代理是否支持开启事务拆分。
	SupportTransactionSplit *bool `json:"support_transaction_split,omitempty"`
}

func (o ProxyInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProxyInfo struct{}"
	}

	return strings.Join([]string{"ProxyInfo", string(data)}, " ")
}

type ProxyInfoSecondsLevelMonitorFunStatus struct {
	value string
}

type ProxyInfoSecondsLevelMonitorFunStatusEnum struct {
	OFF ProxyInfoSecondsLevelMonitorFunStatus
	ON  ProxyInfoSecondsLevelMonitorFunStatus
}

func GetProxyInfoSecondsLevelMonitorFunStatusEnum() ProxyInfoSecondsLevelMonitorFunStatusEnum {
	return ProxyInfoSecondsLevelMonitorFunStatusEnum{
		OFF: ProxyInfoSecondsLevelMonitorFunStatus{
			value: "off",
		},
		ON: ProxyInfoSecondsLevelMonitorFunStatus{
			value: "on",
		},
	}
}

func (c ProxyInfoSecondsLevelMonitorFunStatus) Value() string {
	return c.value
}

func (c ProxyInfoSecondsLevelMonitorFunStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ProxyInfoSecondsLevelMonitorFunStatus) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
