package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListSampleResponse struct {

	// 时间序列对象列表。
	Samples        *[]SampleDataValue `json:"samples,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ListSampleResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSampleResponse struct{}"
	}

	return strings.Join([]string{"ListSampleResponse", string(data)}, " ")
}
