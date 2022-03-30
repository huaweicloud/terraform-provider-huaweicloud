package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListGaussMySqlInstancesResponse struct {
	// 实例列表信息。

	Instances *[]MysqlInstanceListInfo `json:"instances,omitempty"`
	// 总记录数。

	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListGaussMySqlInstancesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGaussMySqlInstancesResponse struct{}"
	}

	return strings.Join([]string{"ListGaussMySqlInstancesResponse", string(data)}, " ")
}
