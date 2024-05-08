package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowNodePoolConfigurationDetailsResponse Response Object
type ShowNodePoolConfigurationDetailsResponse struct {

	// 获取指定节点池配置参数列表返回体
	Body           map[string][]PackageOptions `json:"body,omitempty"`
	HttpStatusCode int                         `json:"-"`
}

func (o ShowNodePoolConfigurationDetailsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowNodePoolConfigurationDetailsResponse struct{}"
	}

	return strings.Join([]string{"ShowNodePoolConfigurationDetailsResponse", string(data)}, " ")
}
