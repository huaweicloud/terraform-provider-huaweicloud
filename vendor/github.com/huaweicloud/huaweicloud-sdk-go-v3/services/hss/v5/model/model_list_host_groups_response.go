package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListHostGroupsResponse Response Object
type ListHostGroupsResponse struct {

	// 总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 服务器组列表
	DataList       *[]HostGroupItem `json:"data_list,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o ListHostGroupsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHostGroupsResponse struct{}"
	}

	return strings.Join([]string{"ListHostGroupsResponse", string(data)}, " ")
}
