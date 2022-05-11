package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 时间序列查询参数详情。
type SeriesApiQueryItemParam struct {

	// 通过该数组传递的参数信息进行时间序列查询。
	Series []QuerySeriesOptionParam `json:"series"`
}

func (o SeriesApiQueryItemParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SeriesApiQueryItemParam struct{}"
	}

	return strings.Join([]string{"SeriesApiQueryItemParam", string(data)}, " ")
}
