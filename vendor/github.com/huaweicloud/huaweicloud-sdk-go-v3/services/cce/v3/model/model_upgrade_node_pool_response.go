package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeNodePoolResponse Response Object
type UpgradeNodePoolResponse struct {

	// API类型，固定值“Job”，该值不可修改。
	Kind *string `json:"kind,omitempty"`

	// API版本，固定值“v3”，该值不可修改。
	ApiVersion *string `json:"apiVersion,omitempty"`

	Metadata *JobMetadata `json:"metadata,omitempty"`

	Spec *JobSpec `json:"spec,omitempty"`

	Status         *JobStatus `json:"status,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o UpgradeNodePoolResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeNodePoolResponse struct{}"
	}

	return strings.Join([]string{"UpgradeNodePoolResponse", string(data)}, " ")
}
