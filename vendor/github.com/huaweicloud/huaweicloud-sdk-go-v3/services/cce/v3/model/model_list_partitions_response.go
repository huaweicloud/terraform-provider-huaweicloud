package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPartitionsResponse Response Object
type ListPartitionsResponse struct {

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	Items          *[]Partition `json:"items,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListPartitionsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPartitionsResponse struct{}"
	}

	return strings.Join([]string{"ListPartitionsResponse", string(data)}, " ")
}
