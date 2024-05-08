package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutoLaunchsResponse Response Object
type ListAutoLaunchsResponse struct {

	// 自启动项总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 自启动项列表
	DataList       *[]AutoLauchResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                      `json:"-"`
}

func (o ListAutoLaunchsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutoLaunchsResponse struct{}"
	}

	return strings.Join([]string{"ListAutoLaunchsResponse", string(data)}, " ")
}
