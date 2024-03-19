package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAssetStatisticRequest Request Object
type ShowAssetStatisticRequest struct {

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// host id
	HostId *string `json:"host_id,omitempty"`

	// 类别，默认为host，包含如下： - host：主机 - container：容器
	Category *string `json:"category,omitempty"`
}

func (o ShowAssetStatisticRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAssetStatisticRequest struct{}"
	}

	return strings.Join([]string{"ShowAssetStatisticRequest", string(data)}, " ")
}
