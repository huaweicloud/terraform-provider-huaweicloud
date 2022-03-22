package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListGaussMySqlErrorLogRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID

	InstanceId string `json:"instance_id"`
	// 开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。

	StartDate string `json:"start_date"`
	// 结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。

	EndDate string `json:"end_date"`
	// 索引位置，偏移量。从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数

	Offset *int32 `json:"offset,omitempty"`
	// 查询记录数。默认为100，不能为负数，最小值为1，最大值为100

	Limit *int32 `json:"limit,omitempty"`
	// 日志级别

	Level *string `json:"level,omitempty"`
	// 节点ID

	NodeId *string `json:"node_id,omitempty"`
}

func (o ListGaussMySqlErrorLogRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGaussMySqlErrorLogRequest struct{}"
	}

	return strings.Join([]string{"ListGaussMySqlErrorLogRequest", string(data)}, " ")
}
