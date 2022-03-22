package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowGaussMySqlQuotasResponse struct {
	// 资源列表对象。

	QuotaList *[]Quota `json:"quota_list,omitempty"`
	// 配额记录的条数。

	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowGaussMySqlQuotasResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlQuotasResponse struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlQuotasResponse", string(data)}, " ")
}
