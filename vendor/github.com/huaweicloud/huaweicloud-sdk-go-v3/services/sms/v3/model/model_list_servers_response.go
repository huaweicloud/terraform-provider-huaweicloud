package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListServersResponse struct {
	// 符合查询条件的源端总数量，不受limit和offset影响

	Count *int32 `json:"count,omitempty"`
	// 批量查询的源端服务器详列表

	SourceServers  *[]SourceServersResponseBody `json:"source_servers,omitempty"`
	HttpStatusCode int                          `json:"-"`
}

func (o ListServersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListServersResponse struct{}"
	}

	return strings.Join([]string{"ListServersResponse", string(data)}, " ")
}
