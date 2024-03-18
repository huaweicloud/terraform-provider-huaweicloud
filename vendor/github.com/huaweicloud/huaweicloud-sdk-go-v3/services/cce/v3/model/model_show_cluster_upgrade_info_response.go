package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClusterUpgradeInfoResponse Response Object
type ShowClusterUpgradeInfoResponse struct {

	// 类型
	Kind *string `json:"kind,omitempty"`

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	Metadata *Metadata `json:"metadata,omitempty"`

	Spec *UpgradeInfoSpec `json:"spec,omitempty"`

	Status         *UpgradeInfoStatus `json:"status,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ShowClusterUpgradeInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterUpgradeInfoResponse struct{}"
	}

	return strings.Join([]string{"ShowClusterUpgradeInfoResponse", string(data)}, " ")
}
