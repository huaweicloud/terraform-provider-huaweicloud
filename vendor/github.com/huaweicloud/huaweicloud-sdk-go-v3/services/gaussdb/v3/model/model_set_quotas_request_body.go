package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SetQuotasRequestBody struct {
	// 资源列表对象。

	QuotaList []SetQuota `json:"quota_list"`
}

func (o SetQuotasRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetQuotasRequestBody struct{}"
	}

	return strings.Join([]string{"SetQuotasRequestBody", string(data)}, " ")
}
