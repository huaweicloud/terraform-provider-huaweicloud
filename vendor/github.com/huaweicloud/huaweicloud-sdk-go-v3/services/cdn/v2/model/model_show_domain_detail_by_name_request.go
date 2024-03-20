package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDomainDetailByNameRequest Request Object
type ShowDomainDetailByNameRequest struct {

	// 加速域名名称。
	DomainName string `json:"domain_name"`

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用帐号调用接口时，该参数必传。  您可以通过调用企业项目管理服务（EPS）的查询企业项目列表接口（ListEnterpriseProject）查询企业项目id。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ShowDomainDetailByNameRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainDetailByNameRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainDetailByNameRequest", string(data)}, " ")
}
