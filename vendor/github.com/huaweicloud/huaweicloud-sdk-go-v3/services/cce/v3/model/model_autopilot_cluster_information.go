package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutopilotClusterInformation
type AutopilotClusterInformation struct {
	Spec *AutopilotClusterInformationSpec `json:"spec"`

	Metadata *AutopilotClusterMetadataForUpdate `json:"metadata,omitempty"`
}

func (o AutopilotClusterInformation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotClusterInformation struct{}"
	}

	return strings.Join([]string{"AutopilotClusterInformation", string(data)}, " ")
}
