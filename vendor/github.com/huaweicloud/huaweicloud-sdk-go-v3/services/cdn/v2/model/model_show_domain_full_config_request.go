package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDomainFullConfigRequest Request Object
type ShowDomainFullConfigRequest struct {

	// 加速域名。
	DomainName string `json:"domain_name"`

	// 企业项目ID， all：所有项目。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 取值为auth_key，用来查询鉴权KEY和鉴权备KEY的值。
	ShowSpecialConfigs *string `json:"show_special_configs,omitempty"`
}

func (o ShowDomainFullConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainFullConfigRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainFullConfigRequest", string(data)}, " ")
}
