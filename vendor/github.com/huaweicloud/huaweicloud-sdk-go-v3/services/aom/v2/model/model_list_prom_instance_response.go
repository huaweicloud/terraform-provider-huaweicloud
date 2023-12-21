package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPromInstanceResponse Response Object
type ListPromInstanceResponse struct {

	// Prometheus实例名称列表。
	Prometheus     *[]PromInstanceEpsModel `json:"prometheus,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o ListPromInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPromInstanceResponse struct{}"
	}

	return strings.Join([]string{"ListPromInstanceResponse", string(data)}, " ")
}
