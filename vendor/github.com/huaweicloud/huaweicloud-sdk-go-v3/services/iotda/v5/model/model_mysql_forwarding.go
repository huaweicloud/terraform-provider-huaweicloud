package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// MysqlForwarding MySql配置信息
type MysqlForwarding struct {
	Address *NetAddress `json:"address"`

	// **参数说明**：连接MYSQL数据库的库名。 **取值范围**：长度不超过64，只允许字母、数字、下划线（_）、连接符（-）的组合。
	DbName string `json:"db_name"`

	// **参数说明**：连接MYSQL数据库的用户名
	Username string `json:"username"`

	// **参数说明**：连接MYSQL数据库的密码
	Password string `json:"password"`

	// **参数说明**：客户端是否使用SSL连接服务端，默认为true
	EnableSsl *bool `json:"enable_ssl,omitempty"`

	// **参数说明**：MYSQL数据库的表名
	TableName string `json:"table_name"`

	// **参数说明**：MYSQL数据库的列和流转数据的对应关系列表。
	ColumnMappings []ColumnMapping `json:"column_mappings"`
}

func (o MysqlForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlForwarding struct{}"
	}

	return strings.Join([]string{"MysqlForwarding", string(data)}, " ")
}
