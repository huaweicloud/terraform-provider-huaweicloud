package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListSeriesResponse struct {

	// 时间序列对象列表。
	Series *[]SeriesQueryItemResult `json:"series,omitempty"`

	MetaData       *MetaDataSeries `json:"meta_data,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ListSeriesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSeriesResponse struct{}"
	}

	return strings.Join([]string{"ListSeriesResponse", string(data)}, " ")
}
