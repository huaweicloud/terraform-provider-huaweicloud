package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDomainRequest Request Object
type ShowDomainRequest struct {

	// 直播域名，如果不设置此字段，则返回租户所有的域名信息
	Domain *string `json:"domain,omitempty"`

	// 企业项目ID，如果不设置此字段，则不进行该字段过滤，返回所有域名信息
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ShowDomainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainRequest", string(data)}, " ")
}
