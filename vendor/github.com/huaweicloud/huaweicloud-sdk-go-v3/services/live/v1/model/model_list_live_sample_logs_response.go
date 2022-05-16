package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListLiveSampleLogsResponse struct {

	// 符合查询条件的总条目数
	Total *int32 `json:"total,omitempty"`

	// 播放域名
	Domain *string `json:"domain,omitempty"`

	// 日志信息列表
	Logs           *[]LogInfo `json:"logs,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o ListLiveSampleLogsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLiveSampleLogsResponse struct{}"
	}

	return strings.Join([]string{"ListLiveSampleLogsResponse", string(data)}, " ")
}
