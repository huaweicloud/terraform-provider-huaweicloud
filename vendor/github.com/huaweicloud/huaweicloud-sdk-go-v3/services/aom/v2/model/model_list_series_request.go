package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListSeriesRequest struct {

	// 用于限制本次返回的结果数据条数。 取值范围(0,1000]，默认值为1000。
	Limit *string `json:"limit,omitempty"`

	// 分页查询起始位置，为非负整数。
	Offset *string `json:"offset,omitempty"`

	Body *SeriesApiQueryItemParam `json:"body,omitempty"`
}

func (o ListSeriesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSeriesRequest struct{}"
	}

	return strings.Join([]string{"ListSeriesRequest", string(data)}, " ")
}
