package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAppChangeHistoriesResponse Response Object
type ListAppChangeHistoriesResponse struct {

	// 软件变动总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 软件历史变动记录列表
	DataList       *[]AppChangeResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                      `json:"-"`
}

func (o ListAppChangeHistoriesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAppChangeHistoriesResponse struct{}"
	}

	return strings.Join([]string{"ListAppChangeHistoriesResponse", string(data)}, " ")
}
