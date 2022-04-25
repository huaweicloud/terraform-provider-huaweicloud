package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListSampleRequest struct {

	// 用于对查询到的时序数据进行断点插值，默认值为-1。 -1：断点处使用-1进行表示。 0 ：断点处使用0进行表示。 null：断点处使用null进行表示。 average：断点处使用前后邻近的有效数据的平均值进行表示，如果不存在有效数据则使用null进行表示。
	FillValue *string `json:"fill_value,omitempty"`

	Body *QuerySampleParam `json:"body,omitempty"`
}

func (o ListSampleRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSampleRequest struct{}"
	}

	return strings.Join([]string{"ListSampleRequest", string(data)}, " ")
}
