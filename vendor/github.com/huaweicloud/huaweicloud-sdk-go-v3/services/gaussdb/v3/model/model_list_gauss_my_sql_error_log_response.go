package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListGaussMySqlErrorLogResponse struct {
	// 错误日志具体信息。

	ErrorLogList *[]MysqlErrorLogList `json:"error_log_list,omitempty"`
	// 总记录数。

	TotalRecord    *int32 `json:"total_record,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListGaussMySqlErrorLogResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGaussMySqlErrorLogResponse struct{}"
	}

	return strings.Join([]string{"ListGaussMySqlErrorLogResponse", string(data)}, " ")
}
