package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HistoryDatabaseInstance PostgreSQL查询可恢复库的实例信息
type HistoryDatabaseInstance struct {

	// 实例ID
	Id *string `json:"id,omitempty"`

	// 实例名称
	Name *string `json:"name,omitempty"`

	// 表的个数
	TotalTables *int32 `json:"total_tables,omitempty"`

	// 数据库信息
	Databases *[]HistoryDatabaseInfo `json:"databases,omitempty"`
}

func (o HistoryDatabaseInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HistoryDatabaseInstance struct{}"
	}

	return strings.Join([]string{"HistoryDatabaseInstance", string(data)}, " ")
}
