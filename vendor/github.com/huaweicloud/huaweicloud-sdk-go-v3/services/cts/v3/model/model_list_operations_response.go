package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListOperationsResponse Response Object
type ListOperationsResponse struct {

	// 全量云服务的操作事件列表。
	Operations     *[]ListOperation `json:"operations,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListOperationsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOperationsResponse struct{}"
	}

	return strings.Join([]string{"ListOperationsResponse", string(data)}, " ")
}
