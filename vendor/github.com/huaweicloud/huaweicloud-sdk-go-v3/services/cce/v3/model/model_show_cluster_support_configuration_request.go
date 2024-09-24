package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClusterSupportConfigurationRequest Request Object
type ShowClusterSupportConfigurationRequest struct {

	// 集群ID，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterId string `json:"cluster_id"`

	// 集群类型，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterType string `json:"cluster_type"`

	// 集群版本，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	ClusterVersion string `json:"cluster_version"`

	// 集群网络类型，获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。
	NetworkMode string `json:"network_mode"`
}

func (o ShowClusterSupportConfigurationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterSupportConfigurationRequest struct{}"
	}

	return strings.Join([]string{"ShowClusterSupportConfigurationRequest", string(data)}, " ")
}
