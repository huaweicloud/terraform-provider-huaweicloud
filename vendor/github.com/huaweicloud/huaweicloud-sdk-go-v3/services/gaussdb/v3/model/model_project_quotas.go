package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProjectQuotas struct {
	// 资源列表对象。

	Resources []Resource `json:"resources"`
}

func (o ProjectQuotas) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProjectQuotas struct{}"
	}

	return strings.Join([]string{"ProjectQuotas", string(data)}, " ")
}
