package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListDomainLogsResponse struct {

	// 日志总数。
	Total *int32 `json:"total,omitempty"`

	// 日志列表数据。
	Logs           *[]CdnLog `json:"logs,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o ListDomainLogsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDomainLogsResponse struct{}"
	}

	return strings.Join([]string{"ListDomainLogsResponse", string(data)}, " ")
}
