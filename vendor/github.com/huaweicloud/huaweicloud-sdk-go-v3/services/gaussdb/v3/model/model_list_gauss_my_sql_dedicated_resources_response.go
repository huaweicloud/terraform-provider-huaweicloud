package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListGaussMySqlDedicatedResourcesResponse struct {
	// 专属资源池信息

	Resources *[]DedicatedResource `json:"resources,omitempty"`
	// 专属资源池数量

	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListGaussMySqlDedicatedResourcesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGaussMySqlDedicatedResourcesResponse struct{}"
	}

	return strings.Join([]string{"ListGaussMySqlDedicatedResourcesResponse", string(data)}, " ")
}
