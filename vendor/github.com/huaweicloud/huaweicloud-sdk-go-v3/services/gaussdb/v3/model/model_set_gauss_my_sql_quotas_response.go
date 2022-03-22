package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type SetGaussMySqlQuotasResponse struct {
	// 资源列表对象。

	QuotaList      *[]SetQuota `json:"quota_list,omitempty"`
	HttpStatusCode int         `json:"-"`
}

func (o SetGaussMySqlQuotasResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetGaussMySqlQuotasResponse struct{}"
	}

	return strings.Join([]string{"SetGaussMySqlQuotasResponse", string(data)}, " ")
}
