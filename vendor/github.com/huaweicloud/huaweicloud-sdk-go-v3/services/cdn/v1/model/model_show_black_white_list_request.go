package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowBlackWhiteListRequest struct {

	// 需要查询IP黑白名单的域名id。获取方法请参见查询加速域名。
	DomainId string `json:"domain_id"`

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ShowBlackWhiteListRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBlackWhiteListRequest struct{}"
	}

	return strings.Join([]string{"ShowBlackWhiteListRequest", string(data)}, " ")
}
