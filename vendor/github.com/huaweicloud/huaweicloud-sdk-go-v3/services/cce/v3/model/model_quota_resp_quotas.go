package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// QuotaRespQuotas 模板配额
type QuotaRespQuotas struct {

	// 资源
	Resources *[]QuotaRespQuotasResources `json:"resources,omitempty"`
}

func (o QuotaRespQuotas) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QuotaRespQuotas struct{}"
	}

	return strings.Join([]string{"QuotaRespQuotas", string(data)}, " ")
}
