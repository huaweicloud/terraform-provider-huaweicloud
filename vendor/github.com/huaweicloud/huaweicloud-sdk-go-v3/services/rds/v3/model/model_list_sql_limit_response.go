package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSqlLimitResponse Response Object
type ListSqlLimitResponse struct {

	// SQL限流总数
	Count *int32 `json:"count,omitempty"`

	// SQL限流详情
	SqlLimitObjects *[]SqlLimitObject `json:"sql_limit_objects,omitempty"`
	HttpStatusCode  int               `json:"-"`
}

func (o ListSqlLimitResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSqlLimitResponse struct{}"
	}

	return strings.Join([]string{"ListSqlLimitResponse", string(data)}, " ")
}
