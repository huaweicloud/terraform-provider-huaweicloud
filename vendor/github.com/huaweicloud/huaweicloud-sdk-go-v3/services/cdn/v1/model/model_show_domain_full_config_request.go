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
}

func (o ShowDomainFullConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainFullConfigRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainFullConfigRequest", string(data)}, " ")
}
