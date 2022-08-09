package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowCertificatesHttpsInfoRequest struct {

	// 每页的数量，取值范围1-10000，不设值时默认值为30。
	PageSize *int32 `json:"page_size,omitempty"`

	// 查询的页码。取值范围1-65535，不设值时默认值为1。
	PageNumber *int32 `json:"page_number,omitempty"`

	// 加速域名。
	DomainName *string `json:"domain_name,omitempty"`

	// 域名所属用户的domain_id。
	UserDomainId *string `json:"user_domain_id,omitempty"`

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ShowCertificatesHttpsInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCertificatesHttpsInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowCertificatesHttpsInfoRequest", string(data)}, " ")
}
