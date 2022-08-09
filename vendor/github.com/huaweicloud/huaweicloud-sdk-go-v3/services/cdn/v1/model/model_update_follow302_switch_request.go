package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateFollow302SwitchRequest struct {

	// 加速域名id。获取方法请参见查询加速域名。
	DomainId string `json:"domain_id"`

	// 当用户开启企业项目功能时，该参数生效，表示修改当前企业项目下加速域名的配置，\"all\"代表所有项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *Follow302StatusRequest `json:"body,omitempty"`
}

func (o UpdateFollow302SwitchRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateFollow302SwitchRequest struct{}"
	}

	return strings.Join([]string{"UpdateFollow302SwitchRequest", string(data)}, " ")
}
