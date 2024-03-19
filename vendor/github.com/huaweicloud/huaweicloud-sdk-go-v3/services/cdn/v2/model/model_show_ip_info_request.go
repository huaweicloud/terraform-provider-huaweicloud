package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowIpInfoRequest Request Object
type ShowIpInfoRequest struct {

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子帐号调用接口时，该参数必传。  您可以通过调用企业项目管理服务（EPS）的查询企业项目列表接口（ListEnterpriseProject）查询企业项目id。
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
