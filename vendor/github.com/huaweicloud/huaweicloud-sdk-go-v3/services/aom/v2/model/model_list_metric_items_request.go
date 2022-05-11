package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListMetricItemsRequest struct {

	// 指标查询方式。
	Type *string `json:"type,omitempty"`

	// 用于限制本次返回的结果数据条数。 取值范围(0,1000]，默认值为1000。
	Limit *string `json:"limit,omitempty"`

	// 分页查询起始位置，为非负整数。
	Start *string `json:"start,omitempty"`

	Body *MetricApiQueryItemParam `json:"body,omitempty"`
}

func (o ListMetricItemsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMetricItemsRequest struct{}"
	}

	return strings.Join([]string{"ListMetricItemsRequest", string(data)}, " ")
}
