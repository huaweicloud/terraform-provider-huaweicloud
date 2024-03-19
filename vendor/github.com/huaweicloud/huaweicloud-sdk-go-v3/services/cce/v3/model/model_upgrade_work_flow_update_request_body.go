package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpgradeWorkFlowUpdateRequestBody struct {
	Status *UpgradeWorkFlowUpdateRequestBodyStatus `json:"status,omitempty"`
}

func (o UpgradeWorkFlowUpdateRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeWorkFlowUpdateRequestBody struct{}"
	}

	return strings.Join([]string{"UpgradeWorkFlowUpdateRequestBody", string(data)}, " ")
}
