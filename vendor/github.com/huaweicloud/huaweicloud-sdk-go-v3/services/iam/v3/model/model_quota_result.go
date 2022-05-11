package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type QuotaResult struct {

	// 资源信息
	Resources *[]Resources `json:"resources,omitempty"`
}

func (o QuotaResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QuotaResult struct{}"
	}

	return strings.Join([]string{"QuotaResult", string(data)}, " ")
}
