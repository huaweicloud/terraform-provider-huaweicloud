package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListImageLocalResponse Response Object
type ListImageLocalResponse struct {

	// 本地镜像总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 本地镜像数据列表
	DataList       *[]ImageLocalInfo `json:"data_list,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o ListImageLocalResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListImageLocalResponse struct{}"
	}

	return strings.Join([]string{"ListImageLocalResponse", string(data)}, " ")
}
