package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResponseCodeSum struct {

	// 1xx请求数
	Sum1xx *int32 `json:"sum1xx,omitempty"`

	// 2xx请求数
	Sum2xx *int32 `json:"sum2xx,omitempty"`

	// 3xx请求数
	Sum3xx *int32 `json:"sum3xx,omitempty"`

	// 4xx请求数
	Sum4xx *int32 `json:"sum4xx,omitempty"`

	// 5xx请求数
	Sum5xx *int32 `json:"sum5xx,omitempty"`
}

func (o ResponseCodeSum) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResponseCodeSum struct{}"
	}

	return strings.Join([]string{"ResponseCodeSum", string(data)}, " ")
}
