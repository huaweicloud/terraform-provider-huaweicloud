package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowAssetStatisticRequest struct {

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// host id
	HostId *string `json:"host_id,omitempty"`
}

func (o ShowAssetStatisticRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAssetStatisticRequest struct{}"
	}

	return strings.Join([]string{"ShowAssetStatisticRequest", string(data)}, " ")
}
