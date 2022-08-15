package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListRangeQueryAomPromGetResponse struct {

	// 响应状态。
	Status *string `json:"status,omitempty"`

	Data           *interface{} `json:"data,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListRangeQueryAomPromGetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRangeQueryAomPromGetResponse struct{}"
	}

	return strings.Join([]string{"ListRangeQueryAomPromGetResponse", string(data)}, " ")
}
