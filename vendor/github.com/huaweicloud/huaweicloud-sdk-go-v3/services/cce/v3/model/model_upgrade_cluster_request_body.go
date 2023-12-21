package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpgradeClusterRequestBody struct {
	Metadata *UpgradeClusterRequestMetadata `json:"metadata"`

	Spec *UpgradeSpec `json:"spec"`
}

func (o UpgradeClusterRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeClusterRequestBody struct{}"
	}

	return strings.Join([]string{"UpgradeClusterRequestBody", string(data)}, " ")
}
