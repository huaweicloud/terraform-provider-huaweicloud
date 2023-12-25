package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpgradeClusterRequestMetadata struct {

	// api版本，默认为v3
	ApiVersion string `json:"apiVersion"`

	// 资源类型，默认为UpgradeTask
	Kind string `json:"kind"`
}

func (o UpgradeClusterRequestMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeClusterRequestMetadata struct{}"
	}

	return strings.Join([]string{"UpgradeClusterRequestMetadata", string(data)}, " ")
}
