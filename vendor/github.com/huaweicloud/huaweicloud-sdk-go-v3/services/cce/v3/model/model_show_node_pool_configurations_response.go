package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowNodePoolConfigurationsResponse Response Object
type ShowNodePoolConfigurationsResponse struct {

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

func (o ShowNodePoolConfigurationsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowNodePoolConfigurationsResponse struct{}"
	}

	return strings.Join([]string{"ShowNodePoolConfigurationsResponse", string(data)}, " ")
}
