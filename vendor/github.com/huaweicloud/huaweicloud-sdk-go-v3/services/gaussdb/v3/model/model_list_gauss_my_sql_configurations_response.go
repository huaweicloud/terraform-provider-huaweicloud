package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListGaussMySqlConfigurationsResponse struct {
	Configurations *[]ConfigurationSummary `json:"configurations,omitempty"`
	// 参数模板的总数。

	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListGaussMySqlConfigurationsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGaussMySqlConfigurationsResponse struct{}"
	}

	return strings.Join([]string{"ListGaussMySqlConfigurationsResponse", string(data)}, " ")
}
