package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateClusterConfigurationsBody 更新指定集群配置参数内容请求体
type UpdateClusterConfigurationsBody struct {

	// API版本，固定值**v3**
	ApiVersion string `json:"apiVersion"`

	// API类型，固定值**Configuration**
	Kind string `json:"kind"`

	Metadata *ConfigurationMetadata `json:"metadata"`

	Spec *ClusterConfigurationsSpec `json:"spec"`
}

func (o UpdateClusterConfigurationsBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateClusterConfigurationsBody struct{}"
	}

	return strings.Join([]string{"UpdateClusterConfigurationsBody", string(data)}, " ")
}
