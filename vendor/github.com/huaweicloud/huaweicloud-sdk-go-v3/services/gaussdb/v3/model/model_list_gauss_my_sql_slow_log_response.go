package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListGaussMySqlSlowLogResponse struct {
	// 错误日志具体信息。

	SlowLogList *[]MysqlSlowLogList `json:"slow_log_list,omitempty"`
	// 慢日志阈值。

	LongQueryTime *string `json:"long_query_time,omitempty"`
	// 总记录数。

	TotalRecord    *int32 `json:"total_record,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListGaussMySqlSlowLogResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGaussMySqlSlowLogResponse struct{}"
	}

	return strings.Join([]string{"ListGaussMySqlSlowLogResponse", string(data)}, " ")
}
