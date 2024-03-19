package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AgentConfig 探针配置
type AgentConfig struct {

	// 探针id
	AgentId *int32 `json:"agent_id,omitempty"`

	// 是否开启数据库影子规则
	DbEnable *bool `json:"db_enable,omitempty"`

	// 数据库影子库名称
	DbShadowRepository *string `json:"db_shadow_repository,omitempty"`

	// 数据库影子库类型
	DbShadowType *string `json:"db_shadow_type,omitempty"`

	// 日志级别，枚举值：INFO，DEBUG，WARN，ERROR
	LogLevel *string `json:"log_level,omitempty"`

	// 日志路径
	LogPath *string `json:"log_path,omitempty"`

	// 影子规则开关
	MainSwitch *bool `json:"main_switch,omitempty"`

	// 是否开启redis影子库规则
	RedisEnable *bool `json:"redis_enable,omitempty"`

	// redis影子库key前缀
	RedisShadowKeyPrefix *string `json:"redis_shadow_key_prefix,omitempty"`

	// redis影子库名称
	RedisShadowRepository *string `json:"redis_shadow_repository,omitempty"`

	// redis影子库类型
	RedisShadowType *string `json:"redis_shadow_type,omitempty"`

	// kafka影子规则开关
	KafkaEnable *bool `json:"kafka_enable,omitempty"`

	// kafka影子topic前缀
	KafkaShadowTopicPrefix *string `json:"kafka_shadow_topic_prefix,omitempty"`

	// 应用日志等级（ALL/TRACE/DEBUG/INFO/WARN/ERROR/OFF）
	AppLogLevel *string `json:"app_log_level,omitempty"`

	// 应用日志路径
	AppLogPath *string `json:"app_log_path,omitempty"`

	// mock规则列表
	MockRuleList *[]MockRuleConfig `json:"mock_rule_list,omitempty"`

	// clickhouse影子规则开关
	ClickhouseEnable *bool `json:"clickhouse_enable,omitempty"`

	// clickhouse影子规则类型
	ClickhouseShadowType *string `json:"clickhouse_shadow_type,omitempty"`

	// clickhouse影子库映射信息
	ClickhouseShadowRepository *string `json:"clickhouse_shadow_repository,omitempty"`

	// 是否开启pulsar影子库规则
	PulsarEnable *bool `json:"pulsar_enable,omitempty"`

	// pulsar影子库前缀
	PulsarShadowTopicPrefix *string `json:"pulsar_shadow_topic_prefix,omitempty"`

	// 自定义配置
	Extensions *string `json:"extensions,omitempty"`
}

func (o AgentConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgentConfig struct{}"
	}

	return strings.Join([]string{"AgentConfig", string(data)}, " ")
}
