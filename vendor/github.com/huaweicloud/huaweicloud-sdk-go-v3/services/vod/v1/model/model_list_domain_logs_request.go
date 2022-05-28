package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListDomainLogsRequest struct {

	// 加速域名，参考格式：www.test1.com。
	DomainName string `json:"domain_name"`

	// 查询开始时间，格式为yyyyMMddHHmmss。 - 查询结果为开始时间之后24小时内的日志数据 - 只能查最近一个月内的数据
	QueryDate string `json:"query_date"`

	// 每页显示日志数量。
	PageSize *int32 `json:"page_size,omitempty"`

	// 当前页数。
	PageNumber *int32 `json:"page_number,omitempty"`
}

func (o ListDomainLogsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDomainLogsRequest struct{}"
	}

	return strings.Join([]string{"ListDomainLogsRequest", string(data)}, " ")
}
