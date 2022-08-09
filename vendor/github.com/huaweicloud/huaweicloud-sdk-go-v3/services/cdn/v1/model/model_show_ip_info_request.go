package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowIpInfoRequest struct {

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// IP地址列表，以“，”分割，最多20个。
	Ips string `json:"ips"`
}

func (o ShowIpInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowIpInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowIpInfoRequest", string(data)}, " ")
}
