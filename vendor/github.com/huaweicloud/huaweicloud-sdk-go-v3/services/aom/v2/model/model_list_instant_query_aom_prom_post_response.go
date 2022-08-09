package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListInstantQueryAomPromPostResponse struct {

	// 响应状态。
	Status *string `json:"status,omitempty"`

	Data           *interface{} `json:"data,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListInstantQueryAomPromPostResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstantQueryAomPromPostResponse struct{}"
	}

	return strings.Join([]string{"ListInstantQueryAomPromPostResponse", string(data)}, " ")
}
