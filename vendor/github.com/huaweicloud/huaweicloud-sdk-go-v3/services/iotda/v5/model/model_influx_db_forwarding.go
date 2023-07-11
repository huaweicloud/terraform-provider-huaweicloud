package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// InfluxDbForwarding InfluxDB配置信息
type InfluxDbForwarding struct {
	Address *NetAddress `json:"address"`

	// **参数说明**：连接InfluxDB数据库的库名,不存在会自动创建
	DbName string `json:"db_name"`

	// **参数说明**：连接InfluxDB数据库的用户名
	Username string `json:"username"`

	// **参数说明**：连接InfluxDB数据库的密码
	Password string `json:"password"`

	// **参数说明**：InfluxDB数据库的measurement,不存在会自动创建
	Measurement string `json:"measurement"`

	// **参数说明**：InfluxDB数据库和流转数据的对应关系列表。
	ColumnMappings []ColumnMapping `json:"column_mappings"`
}

func (o InfluxDbForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InfluxDbForwarding struct{}"
	}

	return strings.Join([]string{"InfluxDbForwarding", string(data)}, " ")
}
