package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AutopilotClusterMetadataForUpdate struct {

	// 集群显示名。  命名规则：以小写字母开头，由小写字母、数字、中划线(-)组成，长度范围4-128位，且不能以中划线(-)结尾。  显示名和其他集群的名称、显示名不可以重复。  为空时表示不进行修改。
	Alias *string `json:"alias,omitempty"`
}

func (o AutopilotClusterMetadataForUpdate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotClusterMetadataForUpdate struct{}"
	}

	return strings.Join([]string{"AutopilotClusterMetadataForUpdate", string(data)}, " ")
}
