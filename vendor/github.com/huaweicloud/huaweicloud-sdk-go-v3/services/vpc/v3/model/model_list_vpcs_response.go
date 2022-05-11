package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListVpcsResponse struct {

	// 请求ID
	RequestId *string `json:"request_id,omitempty"`

	// VPC列表响应体
	Vpcs *[]Vpc `json:"vpcs,omitempty"`

	PageInfo       *PageInfo `json:"page_info,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ListVpcsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVpcsResponse struct{}"
	}

	return strings.Join([]string{"ListVpcsResponse", string(data)}, " ")
}
