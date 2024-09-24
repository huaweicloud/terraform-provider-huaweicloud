package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListContainersResponse Response Object
type ListContainersResponse struct {

	// 容器总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 最近更新时间
	LastUpdateTime *int64 `json:"last_update_time,omitempty"`

	// 容器基本信息列表
	DataList       *[]ContainerBaseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o ListContainersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListContainersResponse struct{}"
	}

	return strings.Join([]string{"ListContainersResponse", string(data)}, " ")
}
