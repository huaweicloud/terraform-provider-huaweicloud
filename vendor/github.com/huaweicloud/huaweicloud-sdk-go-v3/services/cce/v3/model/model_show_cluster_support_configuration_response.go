package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowClusterSupportConfigurationResponse Response Object
type ShowClusterSupportConfigurationResponse struct {

	// 获取指定集群配置项列表返回体
	Body           map[string][]PackageOptions `json:"body,omitempty"`
	HttpStatusCode int                         `json:"-"`
}

func (o ShowClusterSupportConfigurationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowClusterSupportConfigurationResponse struct{}"
	}

	return strings.Join([]string{"ShowClusterSupportConfigurationResponse", string(data)}, " ")
}
