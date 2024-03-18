package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowVulStaticsRequest Request Object
type ShowVulStaticsRequest struct {

	// 企业租户ID，“0”表示默认企业项目，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ShowVulStaticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowVulStaticsRequest struct{}"
	}

	return strings.Join([]string{"ShowVulStaticsRequest", string(data)}, " ")
}
