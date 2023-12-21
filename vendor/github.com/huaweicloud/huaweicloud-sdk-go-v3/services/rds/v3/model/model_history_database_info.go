package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HistoryDatabaseInfo PostgreSQL查询可恢复库的数据库库信息
type HistoryDatabaseInfo struct {

	// 数据库名
	Name *string `json:"name,omitempty"`

	// 表的个数
	TotalTables *int32 `json:"total_tables,omitempty"`
}

func (o HistoryDatabaseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HistoryDatabaseInfo struct{}"
	}

	return strings.Join([]string{"HistoryDatabaseInfo", string(data)}, " ")
}
