package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeVulStatusRequest Request Object
type ChangeVulStatusRequest struct {

	// 缺省值:application/json; charset=utf-8
	ContentType *string `json:"Content-Type,omitempty"`

	// 企业项目ID，“0”表示默认企业项目，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *ChangeVulStatusRequestInfo `json:"body,omitempty"`
}

func (o ChangeVulStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeVulStatusRequest struct{}"
	}

	return strings.Join([]string{"ChangeVulStatusRequest", string(data)}, " ")
}
