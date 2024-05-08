package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListVulHostsResponse Response Object
type ListVulHostsResponse struct {

	// 受影响的云服务器总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 受影响的云服务器列表
	DataList       *[]VulHostInfo `json:"data_list,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ListVulHostsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListVulHostsResponse struct{}"
	}

	return strings.Join([]string{"ListVulHostsResponse", string(data)}, " ")
}
