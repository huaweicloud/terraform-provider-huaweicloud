package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListProcessesHostResponse Response Object
type ListProcessesHostResponse struct {

	// 主机统计信息总数,
	TotalNum *int32 `json:"total_num,omitempty"`

	// 主机统计信息列表
	DataList       *[]ProcessesHostResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                          `json:"-"`
}

func (o ListProcessesHostResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProcessesHostResponse struct{}"
	}

	return strings.Join([]string{"ListProcessesHostResponse", string(data)}, " ")
}
