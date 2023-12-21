package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateNodePoolConfigurationResponse Response Object
type UpdateNodePoolConfigurationResponse struct {

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	// API类型，固定值**Configuration**
	Kind *string `json:"kind,omitempty"`

	Metadata *ConfigurationMetadata `json:"metadata,omitempty"`

	Spec *ClusterConfigurationsSpec `json:"spec,omitempty"`

	// Configuration的状态信息
	Status         *interface{} `json:"status,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o UpdateNodePoolConfigurationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateNodePoolConfigurationResponse struct{}"
	}

	return strings.Join([]string{"UpdateNodePoolConfigurationResponse", string(data)}, " ")
}
